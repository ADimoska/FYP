import sys
import json
import math
import matplotlib.pyplot as plt

def round_up_to_nice(num):
    """Round up to the nearest 'nice' interval based on magnitude (e.g., 100, 500)"""
    if num == 0:
        return 10
    magnitude = 10 ** (len(str(num)) - 1)
    return math.ceil(num / magnitude) * magnitude

def main():
    # Read JSON array from stdin
    data = json.loads(sys.stdin.read())

    if not data:
        print("No data provided", file=sys.stderr)
        return

    max_val = max(data)
    bin_size = round_up_to_nice(max_val // 10 or 1)
    num_bins = (max_val // bin_size) + 1

    # Define bin edges for the histogram
    bin_edges = [i * bin_size for i in range(num_bins + 1)]

    # Plot histogram
    plt.figure(figsize=(10, 6))
    plt.grid(True, alpha=0.3)
    plt.hist(data, bins=bin_edges, edgecolor='dimgrey')

    plt.title("Distribution of Number of Half-Hour Blackouts per House in the Period between 1/07/2012 and 30/06/2013")
    plt.xlabel("Number of Half-Hour Blackouts")
    plt.ylabel("Number of Houses")
    plt.xticks(bin_edges, rotation=45)
    plt.tight_layout()
    plt.savefig("blackout_distribution.png")

if __name__ == "__main__":
    main()
