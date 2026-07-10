package ownerid

import "strconv"

// Parse returns nil if raw is empty (no ownerId supplied — single-user default),
// or a parsed *int64, or an error if raw is non-empty and not a valid integer.
func Parse(raw string) (*int64, error) {
	if raw == "" {
		return nil, nil
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
