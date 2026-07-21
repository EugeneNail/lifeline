package habits

import (
	"github.com/EugeneNail/lifeline/internal/domain"
)

// Icon represents a habit icon.
type Icon int

const (
	// Health
	IconFavorite           Icon = 1
	IconHealthAndSafety    Icon = 2
	IconHealthMetrics      Icon = 3
	IconMedicalServices    Icon = 4
	IconMedicalInformation Icon = 5
	IconMedication         Icon = 6
	IconPill               Icon = 7
	IconVaccines           Icon = 8
	IconDentistry          Icon = 9
	IconStethoscope        Icon = 10
	IconCardiology         Icon = 11
	IconMonitorHeart       Icon = 12
	IconBloodPressure      Icon = 13
	IconThermometer        Icon = 14
	IconBed                Icon = 15
	IconBedtime            Icon = 16
	IconSleep              Icon = 17
	IconSleepScore         Icon = 18
	IconShower             Icon = 19
	IconBathtub            Icon = 20
	IconCleanHands         Icon = 21
	IconSanitizer          Icon = 22
	IconScale              Icon = 23
	IconMonitorWeight      Icon = 24

	// Sports and Movement
	IconFitnessCenter     Icon = 25
	IconExercise          Icon = 26
	IconDirectionsRun     Icon = 27
	IconSprint            Icon = 28
	IconDirectionsWalk    Icon = 29
	IconNordicWalking     Icon = 30
	IconHiking            Icon = 31
	IconDirectionsBike    Icon = 32
	IconPedalBike         Icon = 33
	IconPool              Icon = 34
	IconRowing            Icon = 35
	IconSportsGymnastics  Icon = 36
	IconSportsMartialArts Icon = 37
	IconSportsSoccer      Icon = 38
	IconSportsBasketball  Icon = 39
	IconSportsTennis      Icon = 40
	IconSportsVolleyball  Icon = 41
	IconSportsHandball    Icon = 42
	IconSportsFootball    Icon = 43
	IconSportsRugby       Icon = 44
	IconSportsGolf        Icon = 45
	IconSportsBaseball    Icon = 46
	IconSkateboarding     Icon = 47
	IconSurfing           Icon = 48

	// Food and Drinks
	IconWaterDrop         Icon = 49
	IconWaterBottle       Icon = 50
	IconLocalDrink        Icon = 51
	IconCoffee            Icon = 52
	IconCoffeeMaker       Icon = 53
	IconEmojiFoodBeverage Icon = 54
	IconNutrition         Icon = 55
	IconRestaurant        Icon = 56
	IconRestaurantMenu    Icon = 57
	IconLocalDining       Icon = 58
	IconDinnerDining      Icon = 59
	IconLunchDining       Icon = 60
	IconBreakfastDining   Icon = 61
	IconFastfood          Icon = 62
	IconGrocery           Icon = 63
	IconLocalGroceryStore Icon = 64
	IconCooking           Icon = 65
	IconSkilletCooktop    Icon = 66
	IconBakeryDining      Icon = 67
	IconCake              Icon = 68
	IconEgg               Icon = 69
	IconRiceBowl          Icon = 70
	IconRamenDining       Icon = 71
	IconIcecream          Icon = 72

	// Home and Household
	IconHome                    Icon = 73
	IconHomeAndGarden           Icon = 74
	IconCleaningServices        Icon = 75
	IconCleaning                Icon = 76
	IconCleaningBucket          Icon = 77
	IconMop                     Icon = 78
	IconWash                    Icon = 79
	IconLaundry                 Icon = 80
	IconLocalLaundryService     Icon = 81
	IconDishwasher              Icon = 82
	IconKitchen                 Icon = 83
	IconVacuum                  Icon = 84
	IconIron                    Icon = 85
	IconPottedPlant             Icon = 86
	IconOutdoorGarden           Icon = 87
	IconHomeRepairService       Icon = 88
	IconHomeImprovementAndTools Icon = 89
	IconToolsPowerDrill         Icon = 90
	IconToolsWrench             Icon = 91
	IconToolsLadder             Icon = 92
	IconChair                   Icon = 93
	IconGarage                  Icon = 94
	IconDoorFront               Icon = 95
	IconWindow                  Icon = 96

	// Work and Growth
	IconWork           Icon = 97
	IconLaptop         Icon = 98
	IconComputer       Icon = 99
	IconTerminal       Icon = 100
	IconCode           Icon = 101
	IconCodeBlocks     Icon = 102
	IconDataObject     Icon = 103
	IconBook           Icon = 104
	IconMenuBook       Icon = 105
	IconLibraryBooks   Icon = 106
	IconSchool         Icon = 107
	IconLanguage       Icon = 108
	IconTranslate      Icon = 109
	IconScience        Icon = 110
	IconCalculate      Icon = 111
	IconEditNote       Icon = 112
	IconTaskAlt        Icon = 113
	IconChecklist      Icon = 114
	IconCalendarMonth  Icon = 115
	IconSchedule       Icon = 116
	IconLightbulb      Icon = 117
	IconPsychology     Icon = 118
	IconDraw           Icon = 119
	IconDesignServices Icon = 120

	// Rest and Well-being
	IconMusicNote       Icon = 121
	IconHeadphones      Icon = 122
	IconLibraryMusic    Icon = 123
	IconPiano           Icon = 124
	IconPalette         Icon = 125
	IconBrush           Icon = 126
	IconVideogameAsset  Icon = 127
	IconGames           Icon = 128
	IconMovie           Icon = 129
	IconTheaters        Icon = 130
	IconPhotoCamera     Icon = 131
	IconTravelExplore   Icon = 132
	IconSunny           Icon = 133
	IconSelfImprovement Icon = 134
	IconSpa             Icon = 135
	IconMood            Icon = 136
	IconNature          Icon = 137
	IconForest          Icon = 138
	IconPark            Icon = 139
	IconBeachAccess     Icon = 140
	IconFireplace       Icon = 141
	IconCelebration     Icon = 142
	IconPets            Icon = 143
	IconCamping         Icon = 144
)

var iconNames = map[Icon]string{
	IconFavorite:                "favorite",
	IconHealthAndSafety:         "health_and_safety",
	IconHealthMetrics:           "health_metrics",
	IconMedicalServices:         "medical_services",
	IconMedicalInformation:      "medical_information",
	IconMedication:              "medication",
	IconPill:                    "pill",
	IconVaccines:                "vaccines",
	IconDentistry:               "dentistry",
	IconStethoscope:             "stethoscope",
	IconCardiology:              "cardiology",
	IconMonitorHeart:            "monitor_heart",
	IconBloodPressure:           "blood_pressure",
	IconThermometer:             "thermometer",
	IconBed:                     "bed",
	IconBedtime:                 "bedtime",
	IconSleep:                   "sleep",
	IconSleepScore:              "sleep_score",
	IconShower:                  "shower",
	IconBathtub:                 "bathtub",
	IconCleanHands:              "clean_hands",
	IconSanitizer:               "sanitizer",
	IconScale:                   "scale",
	IconMonitorWeight:           "monitor_weight",
	IconFitnessCenter:           "fitness_center",
	IconExercise:                "exercise",
	IconDirectionsRun:           "directions_run",
	IconSprint:                  "sprint",
	IconDirectionsWalk:          "directions_walk",
	IconNordicWalking:           "nordic_walking",
	IconHiking:                  "hiking",
	IconDirectionsBike:          "directions_bike",
	IconPedalBike:               "pedal_bike",
	IconPool:                    "pool",
	IconRowing:                  "rowing",
	IconSportsGymnastics:        "sports_gymnastics",
	IconSportsMartialArts:       "sports_martial_arts",
	IconSportsSoccer:            "sports_soccer",
	IconSportsBasketball:        "sports_basketball",
	IconSportsTennis:            "sports_tennis",
	IconSportsVolleyball:        "sports_volleyball",
	IconSportsHandball:          "sports_handball",
	IconSportsFootball:          "sports_football",
	IconSportsRugby:             "sports_rugby",
	IconSportsGolf:              "sports_golf",
	IconSportsBaseball:          "sports_baseball",
	IconSkateboarding:           "skateboarding",
	IconSurfing:                 "surfing",
	IconWaterDrop:               "water_drop",
	IconWaterBottle:             "water_bottle",
	IconLocalDrink:              "local_drink",
	IconCoffee:                  "coffee",
	IconCoffeeMaker:             "coffee_maker",
	IconEmojiFoodBeverage:       "emoji_food_beverage",
	IconNutrition:               "nutrition",
	IconRestaurant:              "restaurant",
	IconRestaurantMenu:          "restaurant_menu",
	IconLocalDining:             "local_dining",
	IconDinnerDining:            "dinner_dining",
	IconLunchDining:             "lunch_dining",
	IconBreakfastDining:         "breakfast_dining",
	IconFastfood:                "fastfood",
	IconGrocery:                 "grocery",
	IconLocalGroceryStore:       "local_grocery_store",
	IconCooking:                 "cooking",
	IconSkilletCooktop:          "skillet_cooktop",
	IconBakeryDining:            "bakery_dining",
	IconCake:                    "cake",
	IconEgg:                     "egg",
	IconRiceBowl:                "rice_bowl",
	IconRamenDining:             "ramen_dining",
	IconIcecream:                "icecream",
	IconHome:                    "home",
	IconHomeAndGarden:           "home_and_garden",
	IconCleaningServices:        "cleaning_services",
	IconCleaning:                "cleaning",
	IconCleaningBucket:          "cleaning_bucket",
	IconMop:                     "mop",
	IconWash:                    "wash",
	IconLaundry:                 "laundry",
	IconLocalLaundryService:     "local_laundry_service",
	IconDishwasher:              "dishwasher",
	IconKitchen:                 "kitchen",
	IconVacuum:                  "vacuum",
	IconIron:                    "iron",
	IconPottedPlant:             "potted_plant",
	IconOutdoorGarden:           "outdoor_garden",
	IconHomeRepairService:       "home_repair_service",
	IconHomeImprovementAndTools: "home_improvement_and_tools",
	IconToolsPowerDrill:         "tools_power_drill",
	IconToolsWrench:             "tools_wrench",
	IconToolsLadder:             "tools_ladder",
	IconChair:                   "chair",
	IconGarage:                  "garage",
	IconDoorFront:               "door_front",
	IconWindow:                  "window",
	IconWork:                    "work",
	IconLaptop:                  "laptop",
	IconComputer:                "computer",
	IconTerminal:                "terminal",
	IconCode:                    "code",
	IconCodeBlocks:              "code_blocks",
	IconDataObject:              "data_object",
	IconBook:                    "book",
	IconMenuBook:                "menu_book",
	IconLibraryBooks:            "library_books",
	IconSchool:                  "school",
	IconLanguage:                "language",
	IconTranslate:               "translate",
	IconScience:                 "science",
	IconCalculate:               "calculate",
	IconEditNote:                "edit_note",
	IconTaskAlt:                 "task_alt",
	IconChecklist:               "checklist",
	IconCalendarMonth:           "calendar_month",
	IconSchedule:                "schedule",
	IconLightbulb:               "lightbulb",
	IconPsychology:              "psychology",
	IconDraw:                    "draw",
	IconDesignServices:          "design_services",
	IconMusicNote:               "music_note",
	IconHeadphones:              "headphones",
	IconLibraryMusic:            "library_music",
	IconPiano:                   "piano",
	IconPalette:                 "palette",
	IconBrush:                   "brush",
	IconVideogameAsset:          "videogame_asset",
	IconGames:                   "games",
	IconMovie:                   "movie",
	IconTheaters:                "theaters",
	IconPhotoCamera:             "photo_camera",
	IconTravelExplore:           "travel_explore",
	IconSunny:                   "sunny",
	IconSelfImprovement:         "self_improvement",
	IconSpa:                     "spa",
	IconMood:                    "mood",
	IconNature:                  "nature",
	IconForest:                  "forest",
	IconPark:                    "park",
	IconBeachAccess:             "beach_access",
	IconFireplace:               "fireplace",
	IconCelebration:             "celebration",
	IconPets:                    "pets",
	IconCamping:                 "camping",
}

// NewIcon returns an icon enum value or a violation when the raw value is unsupported.
func NewIcon(rawIcon int) (Icon, domain.Violation) {
	icon := Icon(rawIcon)
	if !icon.IsValid() {
		return 0, domain.NewViolationf("icon must be in range between %d and %d", IconFavorite, IconCamping)
	}

	return icon, nil
}

// IsValid reports whether the icon is one of the supported enum values.
func (icon Icon) IsValid() bool {
	_, ok := iconNames[icon]

	return ok
}

// String returns the Material Symbols ligature name of the icon.
func (icon Icon) String() string {
	name, ok := iconNames[icon]
	if !ok {
		return "Unknown"
	}

	return name
}
