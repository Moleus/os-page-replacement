accesses='7, 6, 4, 12, 20, 3, 8, 11, 13, 6, 1, 8, 24, 11, 6, 10, 14, 7, 13, 8, 23, 8, 22, 21, 11, 18, 22, 23, 22, 12, 21, 6, 17, 18, 20, 19, 9'
# normal (7 frames)
./main -frames 7 -pages 22 -replacer fifo -accesses "$accesses" > results/fifo-7-frames.txt
./main -frames 7 -pages 22 -replacer lru -accesses "$accesses" > results/lru-7-frames.txt
./main -frames 7 -pages 22 -replacer opt -accesses "$accesses" > results/opt-7-frames.txt

# 3 frames
./main -frames 3 -pages 22 -replacer fifo -accesses "$accesses" > results/fifo-3-frames.txt
./main -frames 3 -pages 22 -replacer lru -accesses "$accesses" > results/lru-3-frames.txt
./main -frames 3 -pages 22 -replacer opt -accesses "$accesses" > results/opt-3-frames.txt

# 14 frames
./main -frames 14 -pages 22 -replacer fifo -accesses "$accesses" > results/fifo-14-frames.txt
./main -frames 14 -pages 22 -replacer lru -accesses "$accesses" > results/lru-14-frames.txt
./main -frames 14 -pages 22 -replacer opt -accesses "$accesses" > results/opt-14-frames.txt

# find amount of frames
./main -brute -brute-percent 0.05 -pages 22 -accesses "$accesses" > results/brute-frames-percents.txt
