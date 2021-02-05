#!/bin/bash
dot -Gdpi=300 -Tpng $1.dot -o $1.png
display $1.png
