package mobile

// Mobile interface is for objects that move within the world
type Mobile interface {
	spawn()
	move()
}
