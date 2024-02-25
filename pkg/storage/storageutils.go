package storage

func FindStorage(storages []Storage, id string) Storage {
	for _, s := range storages {
		if s.GetID() == id {
			return s
		}
	}
	return nil
}
