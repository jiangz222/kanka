package kanka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Henry-Sarabia/blank"
)

// Character contains information about a character.
// For more information, visit: https://kanka.io/en-US/docs/1.0/characters
type Character struct {
	SimpleCharacter
	ID             int       `json:"id"`
	ImageFull      string    `json:"image_full"`
	ImageThumb     string    `json:"image_thumb"`
	HasCustomImage bool      `json:"has_custom_image"`
	EntityID       int       `json:"entity_id"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      int       `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      int       `json:"updated_by"`
	Traits         Traits    `json:"traits"`

	Attributes   Attributes   `json:"attributes"`
	EntityEvents EntityEvents `json:"entity_events"`
	EntityFiles  EntityFiles  `json:"entity_files"`
	EntityNotes  EntityNotes  `json:"entity_notes"`
	Relations    Relations    `json:"relations"`
	Inventory    Inventory    `json:"inventory"`
}

// SimpleCharacter contains only the simple information about a character.
// SimpleCharacter is primarily used to create new characters for posting to
// Kanka.
type SimpleCharacter struct {
	Name             string   `json:"name"`
	Entry            string   `json:"entry,omitempty"`
	Title            string   `json:"title,omitempty"`
	Age              string   `json:"age,omitempty"`
	Sex              string   `json:"sex,omitempty"`
	Type             string   `json:"type,omitempty"`
	FamilyID         int      `json:"family_id,omitempty"`
	LocationID       int      `json:"location_id,omitempty"`
	RaceID           int      `json:"race_id,omitempty"`
	Tags             []int    `json:"tags,omitempty"`
	IsDead           bool     `json:"is_dead,omitempty"`
	IsPrivate        bool     `json:"is_private,omitempty"`
	Image            string   `json:"image,omitempty"`
	ImageURL         string   `json:"image_url,omitempty"`
	PersonalityName  []string `json:"personality_name,omitempty"`
	PersonalityEntry []string `json:"personality_entry,omitempty"`
	AppearanceName   []string `json:"appearance_name,omitempty"`
	AppearanceEntry  []string `json:"appearance_entry,omitempty"`
}

// MarshalJSON marshals the SimpleCharacter into its JSON-encoded form if it
// has the required populated fields.
func (sc SimpleCharacter) MarshalJSON() ([]byte, error) {
	if blank.Is(sc.Name) {
		return nil, fmt.Errorf("cannot marshal SimpleCharacter into JSON with a missing Name")
	}

	type alias SimpleCharacter
	return json.Marshal(alias(sc))
}

// Traits wraps a list of character traits.
// Traits exists to satisfy the API's JSON structure.
type Traits struct {
	Data []*Trait `json:"data"`
}

// Trait represents a character's personality or appearance detail.
type Trait struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Entry        string `json:"entry"`
	Section      string `json:"section"`
	IsPrivate    bool   `json:"is_private"`
	DefaultOrder int    `json:"default_order"`
}

// CharacterService handles communication with the Character endpoint.
type CharacterService service

// Index returns the list of all Characters in the Campaign associated with campID.
// If a non-nil time is provided, Index will only return Characters that have
// been changed since that time.
func (cs *CharacterService) Index(campID int, sync *time.Time) ([]*Character, error) {
	end, err := EndpointCampaign.id(campID)
	if err != nil {
		return nil, fmt.Errorf("invalid Campaign ID: %w", err)
	}
	end = end.concat(cs.end)

	if sync != nil {
		end = end.sync(*sync)
	}

	var wrap struct {
		Data []*Character `json:"data"`
	}

	err = cs.client.get(end, &wrap)
	if err != nil {
		return nil, fmt.Errorf("cannot get Character Index from Campaign (ID: %d): %w", campID, err)
	}

	return wrap.Data, nil
}

// Get returns the Character associated with charID from the Campaign
// associated with campID.
func (cs *CharacterService) Get(campID int, charID int) (*Character, error) {
	end, err := EndpointCampaign.id(campID)
	if err != nil {
		return nil, fmt.Errorf("invalid Campaign ID: %w", err)
	}
	end = end.concat(cs.end)

	end, err = end.id(charID)
	if err != nil {
		return nil, fmt.Errorf("invalid Character ID: %w", err)
	}

	var wrap struct {
		Data *Character `json:"data"`
	}

	err = cs.client.get(end, &wrap)
	if err != nil {
		return nil, fmt.Errorf("cannot get Character (ID: %d) from Campaign (ID: %d): %w", charID, campID, err)
	}

	return wrap.Data, nil
}

// Create creates a new Character in the Campaign associated with campID using
// the provided SimpleCharacter data.
// Create returns the newly created Character.
func (cs *CharacterService) Create(campID int, ch SimpleCharacter) (*Character, error) {
	end, err := EndpointCampaign.id(campID)
	if err != nil {
		return nil, fmt.Errorf("invalid Campaign ID: %w", err)
	}
	end = end.concat(cs.end)

	b, err := json.Marshal(ch)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal SimpleCharacter (Name: %s): %w", ch.Name, err)
	}

	var wrap struct {
		Data *Character `json:"data"`
	}

	err = cs.client.post(end, bytes.NewReader(b), &wrap)
	if err != nil {
		return nil, fmt.Errorf("cannot create Character (Name: %s) for Campaign (ID: %d): %w", ch.Name, campID, err)
	}

	return wrap.Data, nil
}

// Update updates an existing Character associated with charID from the
// Campaign associated with campID using the provided SimpleCharacter data.
// Update returns the newly updated Character.
func (cs *CharacterService) Update(campID int, charID int, ch SimpleCharacter) (*Character, error) {
	end, err := EndpointCampaign.id(campID)
	if err != nil {
		return nil, fmt.Errorf("invalid Campaign ID: %w", err)
	}
	end = end.concat(cs.end)

	end, err = end.id(charID)
	if err != nil {
		return nil, fmt.Errorf("invalid Character ID: %w", err)
	}

	b, err := json.Marshal(ch)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal SimpleCharacter (Name: %s): %w", ch.Name, err)
	}

	var wrap struct {
		Data *Character `json:"data"`
	}

	err = cs.client.put(end, bytes.NewReader(b), &wrap)
	if err != nil {
		return nil, fmt.Errorf("cannot update Character (Name: %s) for Campaign (ID: %d): '%w'", ch.Name, campID, err)
	}

	return wrap.Data, nil
}

// Delete deletes an existing Character associated with charID from the
// Campaign associated with campID.
func (cs *CharacterService) Delete(campID int, charID int) error {
	end, err := EndpointCampaign.id(campID)
	if err != nil {
		return fmt.Errorf("invalid Campaign ID: %w", err)
	}
	end = end.concat(cs.end)

	end, err = end.id(charID)
	if err != nil {
		return fmt.Errorf("invalid Character ID: %w", err)
	}

	err = cs.client.delete(end)
	if err != nil {
		return fmt.Errorf("cannot delete Character (ID: %d) for Campaign (ID: %d): %w", charID, campID, err)
	}

	return nil
}
