package game

var (
	ListOfItems = map[string]struct{}{
		"Bird":             {},
		"Fish":             {},
		"Snake":            {},
		"Spider":           {},
		"Butterfly":        {},
		"Hat":              {},
		"Car":              {},
		"Bicycle":          {},
		"Glasses":          {},
		"Chair":            {},
		"Spoon":            {},
		"Fork":             {},
		"Book":             {},
		"Leaf":             {},
		"Tree":             {},
		"Compass":          {},
		"Toothpaste":       {},
		"Fire":             {},
		"Umbrella":         {},
		"Guitar":           {},
		"Candle":           {},
		"Pizza":            {},
		"Cloud":            {},
		"Star":             {},
		"Anchor":           {},
		"Skull":            {},
		"Clock":            {},
		"Dice":             {},
		"Cactus":           {},
		"Camera":           {},
		"Knife":            {},
		"Diamond":          {},
		"Lightbulb":        {},
		"Television":       {},
		"Key":              {},
		"Globe":            {},
		"Ice cream":        {},
		"Crown":            {},
		"Sword":            {},
		"Balloon":          {},
		"Flag":             {},
		"Bow":              {},
		"Arrow":            {},
		"Thermometer":      {},
		"Axe":              {},
		"Kite":             {},
		"Piggy bank":       {},
		"Hourglass":        {},
		"Magnifying glass": {},
		"Battery":          {},
		"Chimney":          {},
		"Rocket":           {},
		"Pistol":           {},
		"Tent":             {},
		"Window":           {},
		"Magnet":           {},
		"Handcuffs":        {},
		"Fountain":         {},
		"Cannon":           {},
		"Cheese":           {},
		"Skeleton":         {},
		"Bat":              {},
		"Teeth":            {},
		"Bridge":           {},
		"Tie":              {},
		"Crab":             {},
		"Volcano":          {},
		"Airplane":         {},
		"House":            {},
		"Piano":            {},
		"Door":             {},
		"Bonfire":          {},
		"Air balloon":      {},
		"Worm":             {},
		"Alien":            {},
		"Snowman":          {},
		"Money":            {},
		"Scissor":          {},
	}
)

func generateWordsInRandomOrder() []string {
	l := len(ListOfItems)
	res := make([]string, l)

	i := 0
	for key := range ListOfItems {
		res[i] = key
		i++
	}

	return res
}
