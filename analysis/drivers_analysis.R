# Drivers Driving Style Analysis in R

# Load required libraries
rLib <- "./analysis/Rlib"
if (!require(tidyverse)) install.packages("tidyverse", dependencies = TRUE, lib = rLib)
library(tidyverse)