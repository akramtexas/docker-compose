package interfaces

// TODO: The abstraction which should be accessed via a global var.
// TODO: Note for (future) refinement / refactoring in that, fßor now, we have this as a duplicate of the (same) abstraction as exists in the edgex-go repo;
// TODO: Eventually, turn this into a shared lib, perhaps...

type ServiceStarter interface {
	Start(service string) error
}

type ServiceStopper interface {
	Stop(service string, ) error
}

type ServiceRestarter interface {
	Restart(service string) error
}
