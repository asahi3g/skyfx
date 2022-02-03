package fps

import (
	"fmt"
	"math"
	"skyfx/utils"
	"time"
)

// TODO : rename *Profile in Profile*

// Globals ...
var g_profiles []Profile
var g_profileStack int32 = 0
var g_profileTree []ProfileId

// ProfileId ...
type ProfileId struct {
	index int32
}

// Profile ...
type Profile struct {
	name   string
	id     ProfileId
	slot   int32
	begin  float64
	delta  float64
	sum    float64
	count  float64
	min    float64
	max    float64
	stack  int32
	offset int32
}

// IsValidProfile ...
func IsValidProfile(this ProfileId) (out bool) {
	out = this.index >= 0 && int(this.index) < len(g_profiles)
	return
}

// InvalidProfile ...
func InvalidProfile() (out ProfileId) {
	out.index = -1
	return
}

// CreateProfile ...
func CreateProfile(name string) (out ProfileId) {

	var profileCount int = len(g_profiles)
	out.index = int32(profileCount)

	for i := 0; i < profileCount; i = i + 1 {
		var profileName string = g_profiles[i].name
		if name == profileName {
			out.index = -1
			i = profileCount
			utils.PanicIf(true, "duplicate entry")
		}
	}

	if out.index >= 0 {
		var profile Profile
		profile.id = out
		profile.name = name
		g_profiles = append(g_profiles, profile)
	}

	return
}

// PrintProfiles ...
func PrintProfiles() {
	for i := 0; i < len(g_profileTree); i = i + 1 {
		PrintProfile(g_profileTree[i])
	}
}

// ClearProfiles ...
func ClearProfiles(full bool) {
	g_profileTree = g_profileTree[:0]

	var profileCount int = len(g_profiles)
	for i := 0; i < profileCount; i = i + 1 {
		var id ProfileId = g_profiles[i].id
		ClearProfile(id, full) // ISSUE : panic when calling fps.ClearProfile(&g_profiles[i], full)
	}
}

// ClearProfile ...
func ClearProfile(this ProfileId, full bool) { // ##1 perf pointer deref vs index
	utils.PanicIfNot(IsValidProfile(this), "invalid id")

	g_profiles[this.index].begin = 0.0
	g_profiles[this.index].delta = 0.0
	g_profiles[this.index].slot = -1
	if full == true {
		g_profiles[this.index].sum = 0.0
		g_profiles[this.index].count = 0.0
		g_profiles[this.index].max = math.MaxFloat64
		g_profiles[this.index].min = math.MaxFloat64
	}
	var stack int32 = g_profiles[this.index].stack
	utils.PanicIf(stack != 0, "invalid stack")
	g_profiles[this.index].stack = 0
}

// PrintProfile ...
func PrintProfile(this ProfileId) {
	utils.PanicIfNot(IsValidProfile(this), "invalid id")
	var cur float64 = g_profiles[this.index].delta
	var count float64 = g_profiles[this.index].count
	var sum float64 = g_profiles[this.index].sum
	var avg float64 = sum / count
	var min float64 = g_profiles[this.index].min
	var max float64 = g_profiles[this.index].max
	var name string = g_profiles[this.index].name
	var offset int32 = g_profiles[this.index].offset
	if count > 0.0 {
		if offset > 0 {
			var i int32 = 0
			for i = 0; i < offset; i = i + 1 {
				fmt.Printf("----")
			}
			fmt.Printf(" ")
		}

		fmt.Printf("%s - #%d, cur %f, avg %f, min %f, max %f\n",
			name,
			int32(count),
			float32(cur*1000.0),
			float32(avg*1000.0),
			float32(min*1000.0),
			float32(max*1000.0))
	}
}

// CreateStartProfile ...
func CreateStartProfile(this ProfileId, name string) (out ProfileId) {
	out = this
	if IsValidProfile(out) == false {
		out = CreateProfile(name)
	}
	StartProfile(out)
	return
}

// StartProfile ...
func StartProfile(this ProfileId) {
	utils.PanicIfNot(IsValidProfile(this), "invalid id")
	g_profiles[this.index].begin = float64(time.Now().UnixNano()) / 1000000000.0 // ##1 refactor
	var stack int32 = g_profiles[this.index].stack
	g_profiles[this.index].stack = stack + 1
	g_profiles[this.index].offset = g_profileStack
	//fmt.Printf("START_PROFILE %d, %s, STACK %d\n", this.index,  g_profiles[this.index].name, g_profileStack)
	g_profileStack = g_profileStack + 1
	var treeIndex int32 = g_profiles[this.index].slot
	if treeIndex <= -1 {
		g_profiles[this.index].slot = int32(len(g_profileTree))
		g_profileTree = append(g_profileTree, this)
	}
}

// StopProfile ...
func StopProfile(this ProfileId) {
	utils.PanicIfNot(IsValidProfile(this), "invalid id")
	var end float64 = float64(time.Now().UnixNano()) / 1000000000.0 // ##1 refactor
	var begin float64 = g_profiles[this.index].begin
	var delta float64 = end - begin
	var sum float64 = g_profiles[this.index].sum
	var count float64 = g_profiles[this.index].count
	var max float64 = g_profiles[this.index].max
	var min float64 = g_profiles[this.index].min
	var stack int32 = g_profiles[this.index].stack
	var previous float64 = g_profiles[this.index].delta
	g_profiles[this.index].delta = previous + delta
	g_profiles[this.index].sum = sum + delta
	g_profiles[this.index].count = count + 1.0
	g_profiles[this.index].max = math.Max(max, delta)
	g_profiles[this.index].min = math.Min(min, delta)
	utils.PanicIf(stack <= 0, "invalid stack")
	g_profiles[this.index].stack = stack - 1
	utils.PanicIf(g_profileStack <= 0, "invalid stack")
	g_profileStack = g_profileStack - 1
	//fmt.Printf("STOP_PROFILE %d, %s, STACK %d\n", this.index, g_profiles[this.index].name, g_profileStack)
}
