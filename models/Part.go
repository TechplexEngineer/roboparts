package models

import (
	"encoding/json"
)

// Part represents an item used in a project. It can represent
// an assembly, single part or COTS(commercial off the shelf) item
// Type ["part", "assembly", "COTS"] //@todo may want to add optional vendor
// ParentPart a pointer to prevent errors with recursive type
type Part struct {
	Common
	PartNumber     string  `json:"part_number"`
	Type           string  `json:"type"`
	Name           string  `json:"name"`
	Notes          string  `json:"notes"`
	Status         string  `json:"status"`
	SourceMaterial string  `json:"source_material"`
	HaveMaterial   bool    `json:"have_material"`
	Quantity       string  `json:"quantity"`
	CutLength      string  `json:"cut_length"`
	Priority       int     `json:"priority"`
	DrawingCreated bool    `json:"drawing_created"`
	ProjectID      uint    `json:"projectId"`
	Project        Project `json:"project"`
	ParentPartId   uint    `json:"parentPartId"`
	ParentPart     *Part   `json:"parent_part"`
	ChildrenParts  []Part  `json:"children_parts" gorm:"foreignKey:ParentPartId"`

	//# Mapping of priority integer stored in database to what is displayed to the user.
	//PRIORITY_MAP = { 0 => "High", 1 => "Normal", 2 => "Low" }

	//# The list of possible part statuses. Key: string stored in database, value: what is displayed to the user.
	//STATUS_MAP = { "designing" => "Design in progress",
	//      "material" => "Material needs to be ordered",
	//      "ordered" => "Waiting for materials",
	//      "drawing" => "Needs drawing",
	//      "ready" => "Ready to manufacture",
	//      "cnc" => "Ready for CNC",
	//      "laser" => "Ready for laser",
	//      "lathe" => "Ready for lathe",
	//      "mill" => "Ready for mill",
	//      "printer" => "Ready for 3D printer",
	//      "router" => "Ready for router",
	//      "manufacturing" => "Manufacturing in progress",
	//      "outsourced" => "Waiting for outsourced manufacturing",
	//      "welding" => "Waiting for welding",
	//      "scotchbrite" => "Waiting for Scotch-Brite",
	//      "anodize" => "Ready for anodize",
	//      "powder" => "Ready for powder coating",
	//      "coating" => "Waiting for coating",
	//      "assembly" => "Waiting for assembly",
	//      "done" => "Done"
	//}

	//many_to_one :project
	//many_to_one :parent_part, :class => self
	//one_to_many :child_parts, :key => :parent_part_id, :class => self

	//# Assigns a part number based on the parent and type and returns a new Part object.
	//def self.generate_number_and_create(project, type, parent_part)
	//parent_part_id = parent_part.nil? ? 0 : parent_part.id
	//parent_part_number = parent_part.nil? ? 0 : parent_part.part_number
	//if type == "part"
	//part_number = Part.filter(:project_id => project.id, :parent_part_id => parent_part_id, :type => "part")
	//.max(:part_number) || parent_part_number
	//part_number += 1
	//else
	//part_number = Part.filter(:project_id => project.id, :type => "assembly").max(:part_number)  || -100
	//part_number += 100
	//end
	//new(:part_number => part_number, :project_id => project.id, :type => type,
	//:parent_part_id => parent_part.nil? ? 0 : parent_part.id)
	//end
	//
	//def full_part_number
	//"#{project.part_number_prefix}-#{type == "assembly" ? "A" : "P"}-%04d" % part_number
	//end
}

func (o Part) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}
