package option

// Database encapsulates information concerning database connections.
type Database struct {
	// ParallparallelIOLimite is the number of concurrent I/O requests your storage subsystem can handle.
	ParallelIOLimit *int

	// SessionBusyRation is the fraction of time that the connection is active executing a statement in the database.
	// If your workload consists of big analytical queries, session_busy_ratio can be up to 1. If your workload consists
	// of many short statements, session_busy_ratio can be close to 0
	SessionBusyRatio *float64

	// AvgParallelism is the average number of backend processes working on a single query.
	AvgParallelism *float64
}

func NewDatabase() *Database {
	return new(Database)
}

// DatabaseWithParallelIOLimit will set the parallelIOLimit field on the database options object.
func DatabaseWithParallelIOLimit(parallelIOLimit int) func(*Database) {
	return func(opts *Database) {
		opts.ParallelIOLimit = &parallelIOLimit
	}
}

// DatabaseWitSessionBusyRatio will set the sessionBusyRatio field on the database options object.
func DatabaseWitSessionBusyRatio(sessionBusyRatio float64) func(*Database) {
	return func(opts *Database) {
		opts.SessionBusyRatio = &sessionBusyRatio
	}
}

// DatabaseWithAvgParallelism will set the avgParallelism field on the database options object.
func DatabaseWithAvgParallelism(avgParallelism float64) func(*Database) {
	return func(opts *Database) {
		opts.AvgParallelism = &avgParallelism
	}
}
