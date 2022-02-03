package gfx

// // Globals ...
// var g_scissors[]math.V4

// // PushScissor ...
// func PushScissor(name string, bounds math.V4) (out math.V4) {
// 	// utils.PanicIfNot(bounds.Z > 0 && bounds.W > 0, "invalid bounds")
// 	out = bounds
// 	var scissorCount int32 = len(g_scissors)
// 	if (scissorCount > 0) {
// 		var previousScissor int32 = scissorCount - 1
// 		var scissor math.V4 = g_scissors[previousScissor]
// 		var x0 float32 = scissor.X
// 		var y0 float32 = scissor.Y
// 		var x1 float32 = x0 + scissor.Z
// 		var y1 float32 = y0 + scissor.W

// 		var x2 float32 = bounds.X
// 		var y2 float32 = bounds.Y
// 		var x3 float32 = x2 + bounds.Z
// 		var y3 float32 = y2 + bounds.W

// 		var maxX float32 = math.Max_f32(x0, x2)
// 		var maxY float32 = math.Max_f32(y0, y2)

// 		var minX float32 = math.Min_f32(x1, x3)
// 		var minY float32 = math.Min_f32(y1, y3)

// 		out.X = maxX
// 		out.Y = maxY
// 		var width float32 = minX - maxX
// 		var height float32 = minY - maxY

// 		out.Z = math.Max_f32(0.0, width)
// 		out.W = math.Max_f32(0.0, height)
// 		var i int32 = 0
// 		//for i = 0; i < scissorCount; i = i + 1 {
// 		//	fmt.Printf("----")
// 		//}
// 		//fmt.Printf("x0 %f, y0 %f, x1 %f, y1 %f, x2 %f, y2 %f, x3 %f, y3 %f\n", x0, y0, x1, y1, x2, y2, x3, y3)
// 		//for i = 0; i < scissorCount; i = i + 1 {
// 		//	fmt.Printf("----")
// 		//}
// 		//fmt.Printf("max %f, %f, min %f, %f\n", maxX, maxY, minX, minY)
// 		//for i = 0; i < scissorCount; i = i + 1 {
// 		//	fmt.Printf("----")
// 		//}
// 		var scissorW float32 = scissor.X + scissor.Z
// 		var scissorH float32 = scissor.Y + scissor.W
// 		//fmt.Printf("%s old scissor %f, %f, %f, %f\n", name, scissor.X, scissor.Y, scissorW, scissorH)
// 		//for i = 0; i < scissorCount; i = i + 1 {
// 		//	fmt.Printf("----")
// 		//}
// //		fmt.Printf("%s bounds  %f, %f, %f, %f\n", name, bounds.X, bounds.Y, bounds.Z, bounds.W)
// //		for i = 0; i < scissorCount; i = i + 1 {
//   //		  fmt.Printf("----")
// 		//}
// 		var outW float32 = out.X + out.Z
// 		var outH float32 = out.Y + out.W
// 		//fmt.Printf("%s new scissor %f, %f, %f, %f\n", name, out.X, out.Y, outW, outH)
// 	} else {
// 		out = bounds
// 		//var i int32 = 0
// 		//for i = 0; i < scissorCount; i = i + 1 {
// 		//	fmt.Printf("----")
// 		//}
// 		//fmt.Printf("%s bounds %f, %f, %f, %f\n", name, bounds.X, bounds.Y, bounds.Z, bounds.W)
// 		//for i = 0; i < scissorCount; i = i + 1 {
// 		//	fmt.Printf("----")
// 		//}
// 		//fmt.Printf("%s out first %f, %f, %f, %f\n", name, out.X, out.Y, out.Z, out.W)

// 	}
// 	g_scissors = append(g_scissors, out)
// 	//if (out.X == 0.0 && out.Y == 0.0 && out.Z == 0 && out.W == 0) {
// 	//	fmt.Printf("%s thostuhsth \n", name)
// 	//	utils.PanicIfNot(false, "")
// 	//}

// 	//var tmp math.V4 =  g_scissors[scissorCount]
// 	//fmt.Printf("%d, tmp readback %f, %f, %f, %f\n", scissorCount, tmp.X, tmp.Y, tmp.Z, tmp.W)

// 	//SetScissor(bounds)
// }

// // PopScissor ...
// func PopScissor() {
// 	var scissorCount int32 = len(g_scissors)
// 	utils.PanicIfNot(scissorCount > 0, "underflow")
// 	g_scissors = resize(g_scissors, scissorCount - 1)
// 	//if (scissorCount > 0) {
// 	//	var scissor int32 = scissorCount - 1
// 		//SetScissor(g_scissors[scissor])
// 	//}
// }
