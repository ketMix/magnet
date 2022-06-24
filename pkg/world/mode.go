package world

// WorldMode represents the type for representing the current game mode
type WorldMode int

// These are our various states.
const (
	PreGameMode  WorldMode = iota // PreGame leads to Build mode.
	BuildMode                     // Build leads to Wave mode.
	WaveMode                      // Wave leads to Wave, Loss, Victory, or PostGame.
	LossMode                      // Loss represents when the core 'splodes. Leads to a restart of the current.
	VictoryMode                   // Victory represents when all waves are finished. Leads to Travel state.
	PostGameMode                  // PostGame is... the final victory...?
)
