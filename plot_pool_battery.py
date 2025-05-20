import sys
import json
import matplotlib.pyplot as plt

def main():
    # Read JSON input from stdin
    data = json.load(sys.stdin)
    data = data[1987280:]

    # Plot the PoolBattery values
    plt.figure(figsize=(10, 6))
    plt.plot(data, marker='o')
    plt.title('Pool Battery Over Time')
    plt.xlabel('Time Step')
    plt.ylabel('Battery Level')
    plt.grid(True)
    plt.tight_layout()
    plt.savefig("pool_battery_plot.png")

if __name__ == "__main__":
    main()
