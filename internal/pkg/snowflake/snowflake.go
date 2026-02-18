package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	epoch          = int64(1609459200000) // 2021-01-01 00:00:00 UTC in milliseconds
	machineIDBits  = uint(10)
	sequenceBits   = uint(12)
	machineIDShift = sequenceBits
	timestampShift = sequenceBits + machineIDBits
	sequenceMask   = int64(-1) ^ (int64(-1) << sequenceBits)
	maxMachineID   = int64(-1) ^ (int64(-1) << machineIDBits)
)

type Generator struct {
	mu        sync.Mutex
	machineID int64
	sequence  int64
	lastStamp int64
}

var (
	instance *Generator
	once     sync.Once
)

func Init(machineID int64) error {
	if machineID < 0 || machineID > maxMachineID {
		return fmt.Errorf("machine ID must be between 0 and %d", maxMachineID)
	}

	once.Do(func() {
		instance = &Generator{
			machineID: machineID,
			sequence:  0,
			lastStamp: -1,
		}
	})

	return nil
}

func GetInstance() *Generator {
	return instance
}

func (g *Generator) Generate() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1e6

	if timestamp < g.lastStamp {
		return 0, fmt.Errorf("clock moved backwards")
	}

	if timestamp == g.lastStamp {
		g.sequence = (g.sequence + 1) & sequenceMask
		if g.sequence == 0 {
			// Sequence overflow, wait for next millisecond
			for timestamp <= g.lastStamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastStamp = timestamp

	id := ((timestamp - epoch) << timestampShift) |
		(g.machineID << machineIDShift) |
		g.sequence

	return id, nil
}

func Generate() (int64, error) {
	if instance == nil {
		return 0, fmt.Errorf("snowflake generator not initialized")
	}
	return instance.Generate()
}
