const express = require('express');
const bodyParser = require('body-parser');

const app = express();
let storage = [];

app.use(bodyParser.json());

app.listen(3000, () => {
    console.log('Server is running on port 3000');
});

app.post('/receipts/process', (req, res) => {
    const receipt = req.body;

    const savedReceipt = {
        id: String(storage.length + 1),
        ...receipt
    };
    storage.push(savedReceipt);
    res.json(savedReceipt);
});

app.get('/receipts/points/:id', (req, res) => {
    const id = req.params.id;
    const receipt = storage.find(r => String(r.id) === String(id));
    if (!receipt) {
        return res.status(404).send('Receipt not found');
    }
    const points = getAllPoints(receipt);
    res.json({ points });
})

app.get('/', (req, res) => {
    res.json(storage);
});

function roundDollarPoints(receipt) {
    let dollarPoints = 0

    if (receipt.total.substring(receipt.total.length - 2) === '00') {
        dollarPoints += 50
    }
    return dollarPoints
}

function getPoints(receipt) {
    let points = 0
    const retailer = receipt.retailer

    for (let i = 0; i < retailer.length; i++) {
        if (retailer.charCodeAt(i) >= 97 && retailer.charCodeAt(i) <= 122) {
            points += 1
        } else if (retailer.charCodeAt(i) >= 65 && retailer.charCodeAt(i) <= 90) {
            points += 1
        } else if (retailer.charCodeAt(i) >= 48 && retailer.charCodeAt(i) <= 57) {
            points += 1
        }
    }
    return points
}

function checkIfMultipleOfFive(receipt) {
    let points = 0

    if (Number(receipt.total % 0.25 === 0)) {
        points += 25
    }
    return points
}

function getPointsFromItems(receipt) {
    let points = 0

    for (let i = 1; i < receipt.items.length; i += 2) {
        if (receipt.items[i])
            points += 5
    }
    return points
}

function pointsForItemLength(receipt) {
    let itemLengthPoints = 0

    for (let i = 0; i < receipt.items.length; i++) {
        if (receipt.items[i].shortDescription.trim().length % 3 === 0) {
            itemLengthPoints += Math.ceil(Number(receipt.items[i].price) * 0.2)
        }
    }
    return itemLengthPoints
}

function pointsForOddPurchaseDate(receipt) {
    let oddPurchaseDatePoints = 0
    const date = receipt.purchaseDate

    if (date.substring(date.length - 1) % 2 !== 0) {
        oddPurchaseDatePoints += 6
    }
    return oddPurchaseDatePoints
}

function pointsForPurchaseTime(receipt) {
    let purchaseTimePoints = 0
    const hours = Number(receipt.purchaseTime.split(':')[0])
    const minutes = Number(receipt.purchaseTime.split(':')[1])

    if (hours === 14 && (minutes > 0 && minutes < 60)) {
        purchaseTimePoints += 10
    } else if (hours === 15) {
        purchaseTimePoints += 10
    }
    return purchaseTimePoints
}


function getAllPoints(receipt) {
    let totalPoints = 0

    totalPoints += getPoints(receipt)
    totalPoints += checkIfMultipleOfFive(receipt)
    totalPoints += getPointsFromItems(receipt)
    totalPoints += pointsForItemLength(receipt)
    totalPoints += pointsForOddPurchaseDate(receipt)
    totalPoints += pointsForPurchaseTime(receipt)
    totalPoints += roundDollarPoints(receipt)

    return totalPoints
}