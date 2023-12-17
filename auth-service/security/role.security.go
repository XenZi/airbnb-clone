package security

import (
	"encoding/json"
	"os"
)

type AccessControl struct {
    PathAccessMap map[string][]string
}

func NewAccessControl() *AccessControl {
    return &AccessControl{
        PathAccessMap: make(map[string][]string),
    }
}

func (ac *AccessControl) LoadAccessConfig(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    err = json.Unmarshal(data, &ac.PathAccessMap)
    if err != nil {
        return err
    }

    return nil
}

func (ac *AccessControl) IsAccessAllowed(userRole string, path string) bool {
    if roles, ok := ac.PathAccessMap[path]; ok {
        for _, role := range roles {
            if role == userRole {
                return true
            }
        }
    }
    return false
}
