# Investigating power trading with Self-Organising Multi Agent Systems codebase

## To run the codes use the command:
```
$ go run .
```

---

## 6.1.1 Experiment: No Common Pool  
Run the code on branch `unlimited-house-battery-1`

---

## 6.1.2 Experiment: Fixed Energy Contribution Threshold  
Navigate to `E2E3auto`. The code is set up for threshold of 0.00 kwh at line 236 in the `main.go`.  
To run the experiment for the rest of the thresholds, **comment out line 236** and **uncomment line 235**.

---

## 6.1.3 Experiment: Dynamic Individual Contribution Thresholds  
This is the only experiment that doesn’t have an own branch.  
There are comments in the `main.go` that provide information how to change the `main.go` file on the branch `E2E3auto` to run these experiments. 

---

## 6.2 Experiment: Community Size  
Navigate to branch `E2E3ESize` and run the code.  
The simulation takes long to run.  
The recorded results are in `Esize.txt`.

---

## 6.3.1 Experiment: Collaboration vs No Collaboration between Two Communities (Unlimited battery)  

- Navigate to branch `E4Mountain2poolsAvg` and run the code to run the experiment for **no collaboration**.  
- Navigate to branch `E4Mountain2poolsAvgCollab` and run the code to run the experiment for **collaboration**. 

---

## 6.3.2 Experiment: Collaboration vs No Collaboration between Two Communities (Limited battery)  

- Navigate to branch `benchmark_mountain` for the **collaboration** experiment  
- Navigate to branch `benchmark_mountain_dont_exchange` for the **no collaboration** experiment

### Setup:
Set `donateTreshold` on **line 142** in the `main.go` to the capacity that you want to test (e.g., `170kwh`, `56kwh`).  
Then on **lines 175 and 176**, replace the placeholder with the capacity value and run the code:

```
pool1 := pool.NewPool(0, <capacity_value>*2*49)
pool2 := pool.NewPool(0, <capacity_value>*2*65)
```

---

## Partial implementation for the future work idea “Favour Tokens”  
Last commit on the branch `main`.
