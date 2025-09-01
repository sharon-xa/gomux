package models

const (
	SECTION_SEPARATOR        = "#" + "\n" + "--------" + "\n"
	SESSION_NAME_VARIABLE    = "SESSION_NAME"
	SESSION_DIR_VARIABLE     = "SESSION_DIR"
	SCRIPT_TEMPLATE_SECTIONS = 10
)

type Direction string

const (
	Horizontal Direction = "h"
	Vertical   Direction = "v"
)

type Layout string

const (
	LayoutTiled          Layout = "tiled"
	LayoutEvenHorizontal Layout = "even-horizontal"
	LayoutEvenVertical   Layout = "even-vertical"
	LayoutMainHorizontal Layout = "main-horizontal"
	LayoutMainVertical   Layout = "main-vertical"
)
