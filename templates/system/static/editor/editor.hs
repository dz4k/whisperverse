init

    set layout to first <.nebula-layout />
    set options to {
        sort: true,
        draggable: ".nebula-layout-item",
        handle: ".nebula-layout-sortable-handle",
        ghostClass: "nebula-layout-sortable-ghost",
        direction: "vertical"
    }
    make a Sortable from layout, options called sortable