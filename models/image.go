package models
import(
	"github.com/satori/go.uuid"
	)

// CREATE TABLE images (id UUID PRIMARY KEY, name string, basic_path string, path_prefix string, folder_path string, 
// 	entity_id string, entity_type string, is_active bool);
type Image struct {
	Base
	Name		string
	BasicURL	string
	ImageURLPrefix	string
	FolderPath	string
	EntityId	uuid.UUID
	EntityType	string
	ImageCategory string
	IsActive	bool
}

