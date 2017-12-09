# Buddhabrot - Draws the fractal of Buddhabrot in the Python progamming language.
# Author:  Joao Nuno Carvalho
# Email:   joaonunocarv@gmail.com
# Date:    2017.12.9
# License: MIT OpenSource License
#
# Description: Draws the Buddhabrot Set.
#              Port in Python based on the C code found on
#              http://paulbourke.net/fractals/buddhabrot/buddha.c
#              See also:
#              http://paulbourke.net/fractals/buddhabrot/
#              https://en.wikipedia.org/wiki/Buddhabrot#Relation_to_the_logistic_map
import multiprocessing

from numba import prange, jit
import numba

import numpy as np
import random
from PIL import Image

import time

# Iterate the Mandelbrot and return TRUE if the point escapes
#@jit(nogil= True, nopython=True, locals={'x': numba.float32, 'y': numba.float32, 'xnew': numba.float32, 'ynew': numba.float32, 'i': numba.int32, 'n': numba.int32})
@jit
def iterate(x0, y0, n, seq_x, seq_y, NMAX):

    x = 0.0
    y = 0.0

    n = 0
    for i in range(0, NMAX):
      xnew = x * x - y * y + x0
      ynew = 2 * x * y + y0
      # seq[i].x = xnew;
      seq_x[i] = xnew
      seq_y[i] = ynew
      if (xnew*xnew + ynew*ynew ) > 10:
          n = 1
          return((True, i))
      x = xnew
      y = ynew
    return((False, -1))


def write_image(filename, image_array, width, height):

    # Find the largest and the lowest density value
    biggest  = np.amax(image_array)
    smallest = np.amin(image_array)

    print("Density value range: " + str(smallest) + " to " + str(biggest))

    # Write the image
    # Raw uncompressed bytes
    im = Image.new("RGB", (width, height))
    pix = im.load()
    for x in range(width):
        for y in range(height):
            ramp = 2*( image_array[x][y] - smallest) / (biggest - smallest)
            if (ramp > 1):
                ramp = 1
            ramp = ramp**0.5
            pix[x,y] = (int(ramp*255), int(ramp*255), int(ramp*255))
    im.save(filename, "PNG")

#@jit(locals={'x': numba.float32, 'y': numba.float32, 'n': numba.int32, 'ix': numba.int32, 'iy': numba.int32, 'i': numba.int32 })
@jit
def buddhabrot(NX, NY, NMAX, TMAX):
    image_out = np.zeros((NX,NY), dtype=np.float32)

    # xy_seq = [(0.0,0.0) for x in range(0, NMAX)]

    xy_seq_x = np.zeros(NMAX, dtype=np.float32)
    xy_seq_y = np.zeros(NMAX, dtype=np.float32)

    n = 0

    NX_2 = NX / 2
    NY_2 = NY / 2

    NX_3 = 0.3 * NX
    NY_3 = 0.3 * NY

    rnd = np.zeros((TMAX * 2), dtype='f')

    #for tt in prange(0, 1000000):
    for tt in range(0, 1000000):
        if (tt%1000 == 0):
            print('iteration ' + str(tt))


        #rnd = np.random.rand(TMAX*2)
        rnd[:] = np.random.rand(*rnd.shape)
        for t in range(0, TMAX):
            # Choose a random  point in same range.
            # x = 6 * random.random() - 3
            # y = 6 * random.random() - 3
            #x = 6 * np.random.rand() - 3
            #y = 6 * np.random.rand() - 3
            x = 6 * rnd[t*2] - 3
            y = 6 * rnd[t*2+1] - 3

            # Determine state of this point, draw if it escapes
            ret_val, n_possible = iterate(x, y, n, xy_seq_x, xy_seq_y, NMAX)
            if (ret_val):
                n = n_possible
                for i in range(0, n):
                    seq_x = xy_seq_x[i]
                    seq_y = xy_seq_y[i]
                    #ix = int(0.3 * NX * (seq_x + 0.5) + NX / 2);
                    #iy = int(0.3 * NY * seq_y + NY / 2);
                    ix = int(NX_3 * (seq_x + 0.5) + NX_2);
                    iy = int(NY_3 * seq_y + NY_2);
                    if ((ix >= 0) and (iy >= 0) and (ix < NX) and (iy < NY)):
                        image_out[iy][ix] += 1

    # Write iamge.
    write_image("buddhabrot_single.png", image_out, NX, NY);

# Image size.
c_NX = 1000
c_NY = 1000

# Lenght of the sequence to test escape status.
# Also known as bailout
c_NMAX = 200

# Number of iterations, multiple of 1 million.
c_TMAX =  10 # 100 # 2000 #1000 # 100


print('CPU_Number: ' + str(multiprocessing.cpu_count()))

start_time = time.time()
buddhabrot(c_NX, c_NY, c_NMAX, c_TMAX)
elapsed_time = time.time() - start_time

print("Elapsed time" + str(elapsed_time))

