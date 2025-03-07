/** 
 * This file adds the "accessibility" extension to htmx, which scans 
 * all content for elements that *should* be focusable, then adds 
 * attributes to guarantees that the browser will focus on them and
 * keyboard event handlers that work as an alternative for "clicks"
 */

htmx.defineExtension("a11y", {

	init: function() {},

	onEvent: function(/** @type {string} */ name, /** @type {Event} */ event) {

		// Only take actions on "htmx:afterProcessNode"
		if (name !== "htmx:afterProcessNode") {
			return;
		}

		// Special rules for links and buttons
		event.target.querySelectorAll("a,[role=link],button,[role=button]").forEach(function(/** @type {HTMLElement} */ node) {
			
			// If tabIndex is not already set, then default it to 0
			if (node.attributes["tabIndex"] == undefined) {
				node.tabIndex = 0
			}

			// If node is focusable (and not already a link or button) then add keyboard handlers for ENTER and SPACE keys
			if (node.tabIndex != -1) {
				if (["A", "BUTTON"].indexOf(node.tagName) == -1) {
					node.addEventListener("keyup", function(event) {
						if ((event.key == "Enter") || (event.key == " ")) {
							htmx.trigger(node, "click")
						}
					})
				}
			}
		})
	}
})
