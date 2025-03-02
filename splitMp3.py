from pydub import AudioSegment
from pydub.silence import split_on_silence
import time
import sys

# Load the MP3 file
if len(sys.argv) != 3:
    print("Usage: python splitMp3.py <input_file> <output_file>")
    sys.exit(1)

audio = AudioSegment.from_mp3(sys.argv[1])

print(f"[{time.strftime('%H:%M:%S')}] Splitting audio on silence Please wait this may take a while...")
# Split based on silence (2 seconds threshold, -40 dB sensitivity)
chunks = split_on_silence(audio, min_silence_len=2000, silence_thresh=-40)

print(f"[{time.strftime('%H:%M:%S')}] Number of chunks: {len(chunks)}")
# Trim silence from each part and repeat 3 times
processed_audio = AudioSegment.empty()    
silence300 = AudioSegment.silent(duration=300) 
silence100 = AudioSegment.silent(duration=100) 

for i, chunk in enumerate(chunks):
    print(f"[{time.strftime('%H:%M:%S')}] Processing chunk {i+1} of {len(chunks)}")
    trimmed_chunk = chunk.strip_silence(silence_len=500, silence_thresh=-40)  # Trim leading/trailing silence
    processed_audio += silence300+(trimmed_chunk + silence100) * 3  # Add silence after each chunk and repeat 3 times

# Export the final processed audio
processed_audio.export(sys.argv[2], format="mp3")

print(f"[{time.strftime('%H:%M:%S')}] Processing complete! Check {sys.argv[2]}")
