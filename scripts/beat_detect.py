#!/usr/bin/env python
"""Beat detection via Librosa. Outputs JSON array of beat timestamps (seconds)."""
import json
import sys


def detect_beats(audio_path: str) -> list[float]:
    import librosa
    y, sr = librosa.load(audio_path, sr=22050)
    tempo, beats = librosa.beat.beat_track(y=y, sr=sr)
    return [float(b * 512 / sr) for b in beats]


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("[]")
        sys.exit(0)
    try:
        result = detect_beats(sys.argv[1])
    except Exception:
        result = []
    print(json.dumps(result))
