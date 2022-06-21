package data

// Wave contains a spawn list and the next wave.
type Wave struct {
	Spawns *SpawnList
	Next   *Wave
}

// SpawnList contains a list of spawns and the next spawns list.
type SpawnList struct {
	Kinds     []string
	Count     int
	Spawnrate int
	Next      *SpawnList
}

// Clone clones our spawn list, wow.
func (sl *SpawnList) Clone() *SpawnList {
	sl2 := &SpawnList{
		Count:     sl.Count,
		Spawnrate: sl.Spawnrate,
	}
	for _, k := range sl.Kinds {
		sl2.Kinds = append(sl2.Kinds, k)
	}
	if sl.Next != nil {
		sl2.Next = sl.Next.Clone()
	}
	return sl2
}

// Clone clones or wave, crazy.
func (w *Wave) Clone() *Wave {
	w2 := &Wave{}

	if w.Spawns != nil {
		w2.Spawns = w.Spawns.Clone()
	}

	if w.Next != nil {
		w2.Next = w.Next.Clone()
	}
	return w2
}
