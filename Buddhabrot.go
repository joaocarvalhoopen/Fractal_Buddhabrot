// Buddhabrot - Draws the fractal of Buddhabrot in the Go progamming language.
// Author:  Joao Nuno Carvalho
// Email:   joaonunocarv@gmail.com
// Date:    2017.12.9
// License: MIT OpenSource License
//
// Description: Draws the Buddhabrot Set.
//              Port in Go based on the C code found on
//              http://paulbourke.net/fractals/buddhabrot/buddha.c
//              See also:
//              http://paulbourke.net/fractals/buddhabrot/
//              https://en.wikipedia.org/wiki/Buddhabrot#Relation_to_the_logistic_map
//
package main

import (
	"fmt"
	"math"
	"image"
	"image/png"
	"image/color"
	"os"
	"math/rand"
	"strconv"
	"time"
//	"flag"
	"log"
	"runtime"
    "runtime/pprof"
//	"sync"
)



// Image size.
const c_NX int32 = 1000
const c_NY int32 = 1000
// Lenght of the sequence to test escape status.
// Also known as bailout
const c_NMAX int32 = 200

// Number of iterations, multiple of 1 million.
const c_TMAX int32 = 1000 //100 // 2000 // 1000 // 100

// NAme of the output file.
const c_filename string = "buddhabrot_single_1000.png"

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")
//var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {

	//flag.Parse()
	cpuprofile := "" // "cpu.prof"
	memprofile := "" // "mem.prof"

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}


	fmt.Println("Starting Buddhabrot.go....")
	// print('CPU_Number: ' + str(multiprocessing.cpu_count()))

	start := time.Now()
	buddhabrot(c_filename, c_NX, c_NY, c_NMAX, c_TMAX)
	t := time.Now()
	elapsed := t.Sub(start)

	fmt.Println("elapsed time: %t ", elapsed)

	fmt.Println("...ending Buddhabrot.go .")



	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}


}

func buddhabrot(filename string, NX, NY, NMAX, TMAX int32) {
	//image_out = np.zeros((NX, NY), dtype = np.float32)
	image_out := [c_NX][c_NY]float64{}

	// xy_seq_x = np.zeros(NMAX, dtype = np.float32)
	// xy_seq_y = np.zeros(NMAX, dtype = np.float32)
	xy_seq_x := [c_NMAX]float64{}
	xy_seq_y := [c_NMAX]float64{}

	var n int32 = 0

	const c_NX_2 int32 = c_NX / 2
	const c_NY_2 int32 = c_NY / 2

	const c_NX_3 float64 = 0.3 * float64(c_NX)
	const c_NY_3 float64 = 0.3 * float64(c_NY)

	for tt := 0; tt < 1000000; tt++ {
		if tt % 1000 == 0 {
			//fmt.Println("iteration " + strconv.Itoa(tt))
			fmt.Println("iteration ", tt)
		}

		for	t := int32(0); t < c_TMAX; t++ {
			// Choose a	random point in	same range.
			x := float64(6*rand.Float32() - 3)
			y := float64(6*rand.Float32() - 3)

			//x := float64(6*rnd_array[t] - 3)
			//y := float64(6*rnd_array[t + TMAX] - 3)

			// Determine state of this point, draw 	if it escapes.
			ret_val, n_possible := iterate(x, y, &xy_seq_x, &xy_seq_y /*, NMAX*/)
			if ret_val {
				n = n_possible
				for	i := int32(0); i < n; i++ {
					seq_x := xy_seq_x[i]
					seq_y := xy_seq_y[i]
					// ix = int(0.3*NX*(seq_x+0.5) + NX/2);
					// iy = int(0.3*NY*seq_y + NY/2);
					ix := int32(float64(c_NX_3)*( seq_x+0.5 ) + float64(c_NX_2))
					iy := int32(float64(c_NY_3)*seq_y + float64(c_NY_2));
					if ((ix >= 0) &&(iy >= 0)) &&
						(ix < c_NX) && (iy < c_NY){
						image_out[iy][ix]++
					}
				}
			}
		}
	}

	// Write
	write_image(filename, &image_out, NX, NY);
}


func iterate(x0, y0 float64,  seq_x, seq_y *[c_NMAX]float64/*, NMAX int32*/) (bool, int32) {

	var x float64 = 0.0
	var y float64 = 0.0

	//var n int32 = 0
	for	i := int32(0); i < c_NMAX; i++{
		xnew := x*x - y*y + x0
		ynew := 2*x*y + y0
		// seq[i].x = xnew;
		seq_x[i] = xnew
		seq_y[i] = ynew
		if (xnew*xnew + ynew*ynew ) > 10 {
			//n = 1
			return true, i
		}
		x = xnew
		y = ynew
	}
	return false, -1
}

func write_image(filename string, image_array *[c_NX][c_NY]float64, width, height int32){

	// Find the largest and the lowest density value
	var biggest float64 = 0
	var smallest float64 = math.MaxInt32

	for _, line := range image_array {
		for _, cell := range line {
			biggest = math.Max(biggest, cell)
			smallest = math.Min(smallest, cell)
		}
	}

	fmt.Println("Density value range: " + strconv.FormatFloat(smallest, 'E', -1, 64) +
		" to " + strconv.FormatFloat(biggest, 'E', -1, 64))

	// Write the image
	// Raw uncompressed bytes

	im := image.NewRGBA(image.Rectangle{image.Point{0,0},image.Point{int(width),int(height)}})
	for x := int32(0); x < width; x++ {
		for y := int32(0); y < width; y++ {
			ramp := 2*( image_array[x][y] - smallest) / (biggest - smallest)
			if ramp > 1{
				ramp = 1
			}
			ramp = math.Pow(ramp, 0.5)
			c := color.RGBA{uint8(255 * ramp),uint8(255 * ramp),uint8(255 * ramp ),255}
			im.Set(int(x), int(y), c)
		}
	}

	myfile, _ := os.Create(filename)

	png.Encode(myfile, im)
}



