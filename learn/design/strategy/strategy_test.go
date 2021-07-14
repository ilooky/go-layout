package strategy

import (
	"fmt"
	"testing"
)

func TestDuck(t *testing.T) {
	mallardDuck := newMallardDuck()

	mallardDuck.display()
	mallardDuck.performFly()
	mallardDuck.performQuack()
	mallardDuck.swim()

	modelDuck := newModelDuck()

	modelDuck.display()
	modelDuck.performFly()
	modelDuck.performQuack()
	modelDuck.swim()

	modelDuck.setFlyingBehaviour(&flyRocketPowered{})
	fmt.Print("New Flying Behaviour of Model Duck: ")
	modelDuck.performFly()
}
