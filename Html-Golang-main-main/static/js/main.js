function addToOrder(name, price) {
    // Implement the logic to add to the order here
    fetch('/add-to-order', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            name: name,
            price: price,
        }),
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to add item to order');
        }
        // Refresh order details after adding an item
        getOrderDetails();
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

function getOrderDetails() {
    fetch('/order-details')
    .then(response => response.json())
    .then(order => {
        // Update the HTML to display order details
        const orderDetailsElement = document.querySelector('.order-details');
        orderDetailsElement.innerHTML = '<div class="food-card">Order Details</div>';
        
        order.items.forEach(item => {
            const foodCardElement = document.createElement('div');
            foodCardElement.classList.add('food-card');
            foodCardElement.textContent = `${item.name} - $${item.price.toFixed(2)}`;
            orderDetailsElement.appendChild(foodCardElement);
        });
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

// Fetch and display initial order details when the page loads
document.addEventListener('DOMContentLoaded', getOrderDetails);