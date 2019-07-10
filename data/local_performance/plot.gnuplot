set terminal png size {{.Width}},{{.Height}} font ",10"
set offset 0,0,graph 0.05, graph 0.05
set xdata time
set timefmt "%Y-%m-%d %H:%M:%S"
set xrange ["{{.FromDateTime}}":]
set yrange [0:]
#set xtics {{.XTics}}
set grid
set key bottom left samplen 0
set datafile separator ","

set bmargin 0
set lmargin 7
set format x ""

set multiplot layout 4,1 rowsfirst
#set ytics 0.5,1
plot "{{.InputFile}}" using 1:2 with lines title "Load Average" lt 1 lw 2
set tmargin 0
set ytics 500,1000
plot "{{.InputFile}}" using 1:3 with lines title "Used RAM, MBytes" lt 3 lw 2
set format x "%H:%M"
set ytics 50,100
plot "{{.InputFile}}" using 1:4 with lines title "Process count" lt 4 lw 2
