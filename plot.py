import json
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
import matplotlib.animation as animation


# Function to load and extract data from JSON
def load_and_extract_data(filepath):
    with open(filepath, 'r') as file:
        data = json.load(file)

    snake_x, snake_y, snake_z = [], [], []
    for snake in data['snakes']:
        for coord in snake['geometry']:
            snake_x.append(coord[0])
            snake_y.append(coord[1])
            snake_z.append(coord[2])

    enemy_x, enemy_y, enemy_z = [], [], []
    for enemy in data['enemies']:
        for coord in enemy['geometry']:
            enemy_x.append(coord[0])
            enemy_y.append(coord[1])
            enemy_z.append(coord[2])

    food_x, food_y, food_z = [], [], []
    for food in data['food']:
        food_x.append(food['c'][0])
        food_y.append(food['c'][1])
        food_z.append(food['c'][2])

    special_x, special_y, special_z = [], [], []
    for coord in data['specialFood']['golden']:
        special_x.append(coord[0])
        special_y.append(coord[1])
        special_z.append(coord[2])

    for coord in data['specialFood']['suspicious']:
        special_x.append(coord[0])
        special_y.append(coord[1])
        special_z.append(coord[2])

    fence_x, fence_y, fence_z = [], [], []
    for coord in data['fences']:
        fence_x.append(coord[0])
        fence_y.append(coord[1])
        fence_z.append(coord[2])

    points = data['points']
    return {
        'snake': (snake_x, snake_y, snake_z),
        'enemy': (enemy_x, enemy_y, enemy_z),
        'food': (food_x, food_y, food_z),
        'special': (special_x, special_y, special_z),
        'fence': (fence_x, fence_y, fence_z),
        'points': points
    }

# Initialize the 3D plot
fig = plt.figure()
ax = fig.add_subplot(111, projection='3d')

# Function to update the plot in real-time
def update(frame):
    ax.clear()
    filepath = 'plot.json'
    data = load_and_extract_data(filepath)

    snake_x, snake_y, snake_z = data['snake']
    enemy_x, enemy_y, enemy_z = data['enemy']
    food_x, food_y, food_z = data['food']
    special_x, special_y, special_z = data['special']
    fence_x, fence_y, fence_z = data['fence']

    ax.scatter(snake_x, snake_y, snake_z, c='blue', marker='o', label='Snakes')
    ax.scatter(food_x, food_y, food_z, c='green', marker='s', label='Food')
    ax.scatter(special_x, special_y, special_z, c='gold', marker='x', label='Special Food')
    ax.scatter(fence_x, fence_y, fence_z, c='black', marker='d', label='Fences')
    ax.scatter(enemy_x, enemy_y, enemy_z, c='red', marker='^', label='Enemy')

    ax.set_xlabel('X Axis')
    ax.set_ylabel('Y Axis')
    ax.set_zlabel('Z Axis')
    points = data['points']
    plt.title(f"3D Visualization of Points - {points} points")
    ax.legend()

# Assign the animation to a variable
ani = animation.FuncAnimation(fig, update, interval=1000)

# Display the plot
plt.show()

