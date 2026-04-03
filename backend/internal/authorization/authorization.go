package authorization

type ScopeSet struct {
	Any bool
	Own bool
}

func BuildScopeSet(scopes []string) ScopeSet {
	set := ScopeSet{}
	for _, s := range scopes {
		switch s {
		case "any":
			set.Any = true
		case "own":
			set.Own = true
		}
	}

	return set
}

func (s ScopeSet) Allows(actorID, ownerID int32) bool {
	if s.Any {
		return true
	}

	return s.Own && actorID == ownerID
}
