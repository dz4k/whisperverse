behavior wysiwyg(name)

	-- WYSIWYG setup
	init 
		-- save links to important DOM nodes
		set element form to closest <form />
		set element input to form.elements[name]
		set element toolbar to first <.wysiwyg-toolbar /> in me
		set element editor to first <.wysiwyg-editor /> in me

		-- configure related DOM nodes
		add [@type="button"] to <button/> in me
		add [@tabIndex=0] to element editor
		add [@contentEditable=true] to element editor

	-- Clicking a toolbar button triggers a command on the content
	on click from <button />
		if closest <.wysiwyg /> to target is not me then
			exit
		end
		set command to target's [@data-command]
		if command is not null then 
			set value to target's [@data-command-value]
			call document.execCommand(command, false, value)
		end

	-- Show the toolbar when focused
	on focus(target) from <.wysiwyg-editor />
		if I am closest <.wysiwyg /> to target
			remove [@hidden] from element toolbar
		end

	-- Hide the toolbar when blured
	on focusout from <.wysiwyg-editor />
		if I am closest <.wysiwyg /> to target
			wait 200ms
			if (<:focus/> in me) is empty then
				add [@hidden=true] to element toolbar
			end
		end

	-- Autosave the WYSIWYG after 15s of inactivity
	on input debounced at 15s
		set element input's value to element editor's innerHTML
		send updated to form
	
	-- Autosave the WYSIWYG whenever it loses focus
	on blur from <.wysiwyg-editor />
		set element input's value to element editor's innerHTML
		send updated to form
