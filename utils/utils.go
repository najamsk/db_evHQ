package utils

//Strings struct
type Strings struct {
}

//ContainsSlice find if key exist inside map
func (s *Strings) ContainsSlice(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

//ContainsSlice find if key exist inside map
func (s *Strings) ContainsMap(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
