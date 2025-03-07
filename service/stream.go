package service

import (
	"time"

	"github.com/benpate/data"
	"github.com/benpate/data/option"
	"github.com/benpate/datatype"
	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/benpate/form"
	"github.com/benpate/schema"
	"github.com/whisperverse/whisperverse/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Stream manages all interactions with the Stream collection
type Stream struct {
	collection            data.Collection
	attachmentService     *Attachment
	templateService       *Template
	draftService          *StreamDraft
	formLibrary           *form.Library
	templateUpdateChannel chan string
	streamUpdateChannel   chan model.Stream
}

// NewStream returns a fully populated Stream service.
func NewStream(collection data.Collection, templateService *Template, draftService *StreamDraft, attachmentService *Attachment, formLibrary *form.Library, templateUpdateChannel chan string, streamUpdateChannel chan model.Stream) Stream {

	return Stream{
		collection:            collection,
		templateService:       templateService,
		draftService:          draftService,
		attachmentService:     attachmentService,
		formLibrary:           formLibrary,
		templateUpdateChannel: templateUpdateChannel,
		streamUpdateChannel:   streamUpdateChannel,
	}
}

/*******************************************
 * REAL-TIME UPDATES
 *******************************************/

// start begins the background watchers used by the Stream Service
func (service *Stream) Watch() {
	for {
		templateID := <-service.templateUpdateChannel
		service.updateStreamsByTemplate(templateID)
	}
}

/*******************************************
 * COMMON DATA FUNCTIONS
 *******************************************/

// List returns an iterator containing all of the Streams who match the provided criteria
func (service *Stream) List(criteria exp.Expression, options ...option.Option) (data.Iterator, error) {
	return service.collection.List(notDeleted(criteria), options...)
}

// Load retrieves an Stream from the database
func (service *Stream) Load(criteria exp.Expression, stream *model.Stream) error {

	if err := service.collection.Load(notDeleted(criteria), stream); err != nil {
		return derp.Wrap(err, "service.Stream", "Error loading Stream", criteria)
	}

	return nil
}

// Save adds/updates an Stream in the database
func (service *Stream) Save(stream *model.Stream, note string) error {

	if err := service.collection.Save(stream, note); err != nil {
		return derp.Wrap(err, "service.Stream", "Error saving Stream", stream, note)
	}

	// NON-BLOCKING: Notify other processes on this server that the stream has been updated
	go func() {
		service.streamUpdateChannel <- *stream
		// fmt.Println("streamService.Save: sent update update to stream: " + stream.Label)
	}()

	// One milisecond delay prevents overlapping stream.CreateDates.  Deal with it.
	// TODO: There has to be a better way than this...
	time.Sleep(1 * time.Millisecond)

	return nil
}

// Delete removes an Stream from the database (virtual delete)
func (service *Stream) Delete(stream *model.Stream, note string) error {

	// Delete this Stream
	if err := service.collection.Delete(stream, note); err != nil {
		return derp.Wrap(err, "service.Stream.Delete", "Error deleting Stream", stream, note)
	}

	go func() {

		// Delete all related Children
		if err := service.DeleteChildren(stream, note); err != nil {
			derp.Report(derp.Wrap(err, "service.Stream.Delete", "Error deleting child streams", stream, note))
		}

		// Delete all related Attachments
		if err := service.attachmentService.DeleteByStream(stream.StreamID, note); err != nil {
			derp.Report(derp.Wrap(err, "service.Stream.Delete", "Error deleting attachments", stream, note))
		}

		// Delete all related Drafts
		if err := service.draftService.Delete(stream, note); err != nil {
			derp.Report(derp.Wrap(err, "service.Stream.Delete", "Error deleting drafts", stream, note))
		}
	}()

	return nil
}

/*******************************************
 * GENERIC DATA FUNCTIONS
 *******************************************/

// New returns a fully initialized model.Stream as a data.Object.
func (service *Stream) ObjectNew() data.Object {
	result := model.NewStream()
	return &result
}

func (service *Stream) ObjectList(criteria exp.Expression, options ...option.Option) (data.Iterator, error) {
	return service.List(criteria, options...)
}

func (service *Stream) ObjectLoad(criteria exp.Expression) (data.Object, error) {
	result := model.NewStream()
	err := service.Load(criteria, &result)
	return &result, err
}

func (service *Stream) ObjectSave(object data.Object, comment string) error {
	return service.Save(object.(*model.Stream), comment)
}

func (service *Stream) ObjectDelete(object data.Object, comment string) error {
	return service.Delete(object.(*model.Stream), comment)
}

func (service *Stream) Debug() datatype.Map {
	return datatype.Map{
		"service": "Stream",
	}
}

/*******************************************
 * CUSTOM QUERIES
 *******************************************/

// DeleteChildren removes all child streams from the provided stream (virtual delete)
func (service *Stream) DeleteChildren(stream *model.Stream, note string) error {

	var child model.Stream
	it, err := service.ListByParent(stream.StreamID)

	if err != nil {
		return derp.Wrap(err, "service.Stream.Delete", "Error listing child streams", stream)
	}

	for it.Next(&child) {
		if err := service.Delete(&child, note); err != nil {
			return derp.Wrap(err, "service.Stream.Delete", "Error deleting child stream", child)
		}
	}

	return nil
}

// ListByParent returns all Streams that match a particular parentID
func (service *Stream) ListByParent(parentID primitive.ObjectID) (data.Iterator, error) {
	return service.List(exp.Equal("parentId", parentID))
}

// ListTopLevel returns all Streams of type FOLDER at the top of the hierarchy
func (service *Stream) ListTopLevel() (data.Iterator, error) {
	return service.List(
		exp.Equal("parentId", primitive.NilObjectID),
		option.SortAsc("rank"),
	)
}

// ListAncestors returns all Streams that are ancestors of the provided stream.
func (service *Stream) ListAncestors(stream *model.Stream) ([]model.Stream, error) {

	result := make([]model.Stream, len(stream.ParentIDs))
	it, err := service.List(exp.In("_id", stream.ParentIDs))

	if err != nil {
		return result, derp.Wrap(err, "service.Stream.ListAncestors", "Error accessing database", stream)
	}

	temp := model.NewStream()

	for it.Next(&temp) {
		result[len(temp.ParentIDs)] = temp
		temp = model.NewStream()
	}

	return result, nil
}

// ListByTemplate returns all Streams that use a particular Template
func (service *Stream) ListByTemplate(template string) (data.Iterator, error) {
	return service.List(exp.Equal("templateId", template))
}

// LoadByToken returns a single Stream that matches a particular Token
func (service *Stream) LoadByToken(token string, result *model.Stream) error {

	// If the token looks like an ObjectID, then try Load by ID first.
	if streamID, err := primitive.ObjectIDFromHex(token); err == nil {
		if err := service.LoadByID(streamID, result); err == nil {
			return nil
		}
	}

	// Default to Load by Token
	return service.Load(exp.Equal("token", token), result)
}

// LoadByID returns a single Stream that matches a particular StreamID
func (service *Stream) LoadByID(streamID primitive.ObjectID, result *model.Stream) error {
	return service.Load(exp.Equal("_id", streamID), result)
}

// LoadBySource locates a single stream that matches the provided SourceURL
func (service *Stream) LoadBySource(parentStreamID primitive.ObjectID, sourceURL string, result *model.Stream) error {

	criteria := exp.
		Equal("parentId", parentStreamID).
		AndEqual("sourceUrl", sourceURL)

	return service.Load(criteria, result)
}

// LoadParent returns the Stream that is the parent of the provided Stream
func (service *Stream) LoadParent(stream *model.Stream, parent *model.Stream) error {

	if !stream.HasParent() {
		return derp.New(404, "service.Stream.LoadParent", "Stream does not have a parent")
	}

	if err := service.LoadByID(stream.ParentID, parent); err != nil {
		return derp.Wrap(err, "service.stream.LoadParent", "Error loading parent", stream)
	}

	return nil
}

// LoadTopLevelByID locates a single stream in the top level of the site hierarchy
func (service *Stream) LoadTopLevelByID(streamID primitive.ObjectID, result *model.Stream) error {

	criteria := exp.
		Equal("_id", streamID).
		AndEqual("parentId", primitive.NilObjectID)

	return service.Load(criteria, result)
}

// ChildTemplates returns an iterator of Templates that can be added as a sub-stream
func (service *Stream) ChildTemplates(stream *model.Stream) []model.Option {
	return service.templateService.ListByContainer(stream.TemplateID)
}

/*******************************************
 * CUSTOM ACTIONS
 *******************************************/

// NewTopLevel creates a new stream at the top level of the tree
func (service *Stream) NewTopLevel(templateID string) (model.Stream, *model.Template, error) {

	template, err := service.templateService.Load(templateID)

	if err != nil {
		return model.Stream{}, nil, derp.Wrap(err, "service.Stream.NewTopLevel", "Cannot find template")
	}

	if !template.CanBeContainedBy("top") {
		return model.Stream{}, template, derp.New(derp.CodeInternalError, "service.Stream.NewTopLevel", "Template cannot be placed at top level", templateID)
	}

	result := model.NewStream()
	result.ParentID = primitive.NilObjectID
	result.ParentIDs = make([]primitive.ObjectID, 0)
	result.TemplateID = templateID

	// TODO: sort order?
	// TODO: presets defined by templates?

	return result, template, nil
}

// NewTopLevel creates a new stream at the top level of the tree
func (service *Stream) NewChild(parent *model.Stream, templateID string) (model.Stream, *model.Template, error) {

	template, err := service.templateService.Load(templateID)

	if err != nil {
		return model.Stream{}, nil, derp.Wrap(err, "service.Stream.NewTopLevel", "Cannot find template")
	}

	if !template.CanBeContainedBy(parent.TemplateID) {
		return model.Stream{}, nil, derp.New(derp.CodeInternalError, "service.Stream.NewTopLevel", "Template cannot be placed at top level", templateID)
	}

	result := model.NewStream()
	result.ParentID = parent.StreamID
	result.ParentIDs = append(parent.ParentIDs, parent.StreamID)
	result.TemplateID = templateID

	// TODO: sort order?
	// TODO: presets defined by templates?

	return result, template, nil
}

// NewTopLevel creates a new stream at the top level of the tree
func (service *Stream) NewSibling(sibling *model.Stream, templateID string) (model.Stream, *model.Template, error) {

	if sibling.HasParent() {
		var parent model.Stream
		if err := service.LoadParent(sibling, &parent); err != nil {
			return model.Stream{}, nil, derp.Wrap(err, "service.Stream.NewSiblling", "Error loading parent Stream")
		}

		return service.NewChild(&parent, templateID)
	}

	return service.NewTopLevel(templateID)
}

// Template returns the Template associated with this Stream
func (service *Stream) Template(templateID string) (*model.Template, error) {
	template, err := service.templateService.Load(templateID)
	return template, err
}

// State returns the detailed State information associated with this Stream
func (service *Stream) State(stream *model.Stream) (model.State, error) {
	return service.templateService.State(stream.TemplateID, stream.StateID)
}

// Schema returns the Schema associated with this Stream
func (service *Stream) Schema(stream *model.Stream) (schema.Schema, error) {
	return service.templateService.Schema(stream.TemplateID)
}

// Action returns the action definition that matches the stream and type provided
func (service *Stream) Action(stream *model.Stream, actionID string) (*model.Action, error) {
	return service.templateService.Action(stream.TemplateID, actionID)
}

// updateStreamsByTemplate pushes every stream that uses a particular template into the streamUpdateChannel.
func (service *Stream) updateStreamsByTemplate(templateID string) {

	iterator, err := service.ListByTemplate(templateID)

	if err != nil {
		derp.Report(derp.Wrap(err, "service.Realtime", "Error Listing Streams for Template", templateID))
		return
	}

	stream := new(model.Stream)

	for iterator.Next(stream) {
		service.streamUpdateChannel <- *stream
		stream = new(model.Stream)
	}
}

// CreatePersonalStream generates a hidden stream that is tightly linked to a specific user
func (service *Stream) CreatePersonalStream(user *model.User, templateID string) (primitive.ObjectID, error) {

	stream := model.NewStream()
	stream.TemplateID = templateID
	stream.ParentID = user.UserID
	stream.AuthorID = user.UserID
	stream.Criteria = model.NewCriteria()
	stream.Criteria.OwnerID = user.UserID

	err := service.Save(&stream, "auto: create inbox")

	return stream.StreamID, err
}
