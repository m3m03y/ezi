import matplotlib.pyplot as plt
import networkx as nx
import pandas as pd
import numpy as np 
from numpy import genfromtxt

def show_graph_with_labels(adjacency_matrix, mylabels):
    rows, cols = np.where(adjacency_matrix == 1)
    edges = zip(rows.tolist(), cols.tolist())
    gr = nx.DiGraph(directed=True)
    all_rows = range(0, adjacency_matrix.shape[0])
    for n in all_rows:
        gr.add_node(n)
    
    gr.add_edges_from(edges)
    nx.draw(gr, node_size=1500, labels=mylabels, with_labels=True)
    plt.show()


# Draw graph
input_data = genfromtxt('adj_matrix.txt', delimiter=',')
adjacency_matrix = input_data[1:,1:]

labels = input_data[0][1:]
labels_dic = {}
for i in range(len(adjacency_matrix)):
    labels_dic[i] = labels[i]
show_graph_with_labels(adjacency_matrix, labels_dic)

# Draw graph with prestige - ranking 1
# input_data = pd.read_csv('ranking1_prestige_labels.txt', index_col=0)
# labels = input_data.columns
# labels_dic = {}
# for i in range(len(adjacency_matrix)):
#     labels_dic[i] = labels[i]
# show_graph_with_labels(adjacency_matrix, labels_dic)

# Draw graph with prestige - ranking 2
# input_data = pd.read_csv('ranking2_prestige_labels.txt', index_col=0)
# labels = input_data.columns
# labels_dic = {}
# for i in range(len(adjacency_matrix)):
#     labels_dic[i] = labels[i]
# show_graph_with_labels(adjacency_matrix, labels_dic)