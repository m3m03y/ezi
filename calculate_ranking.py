import numpy as np
from numpy import genfromtxt
import json

input_data = genfromtxt('adj_matrix.txt', delimiter=',')
adjacency_matrix = input_data[1:,1:]

site_ranking = {}
eigenvalues = np.linalg.eig(adjacency_matrix)

file = open("site_list.txt", "r")
site_names = file.read().split("\n")
site_names.remove('')
labels = input_data[0][1:]

for i in range(len(labels)):
    current = eigenvalues[0][i]
    site_ranking[site_names[i]] = abs(float(current.real))


site_ranking = dict(sorted(site_ranking.items(), key=lambda item: item[1], reverse=True))

prestige_labels_txt = ""
for i in range(len(site_names)):
    value = round(site_ranking.get(site_names[i]),2)
    prestige_labels_txt += f',{i} {value}'

with open('ranking3_prestige_labels.txt', 'w') as fp:
    fp.write(prestige_labels_txt)
    
with open('ranking_3.json', 'w') as fp:
    json.dump(site_ranking, fp)
