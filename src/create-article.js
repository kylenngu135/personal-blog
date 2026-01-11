async function publish() {
    const title_name = document.getElementById('title-input').value;
    const pub_date = document.getElementById('pub-input').value;
    const cont = document.getElementById('content-input').value;

    // check for empty fields
    if (!title_name) {
        alert('Please enter a title.');
        return
    } else if (!pub_date) {
        alert('Please enter a publishing date');
        return
    } else if (!cont) {
        alert('Can\'t publish an article with no content');
        return
    }

    try {
        const response = await fetch('http://localhost:8080/publish', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title: title_name,
                date: pub_date,
                content: cont
            })
        });

        const data = await response.json();

        if (data.error) {
            alert('Error: ' + data.error);
        } else {
            alert('SUCCESS!')
        }
    } catch (error) {
        alert('Failed to connect to server: ' + error.message)
    }
}

document.addEventListener('DOMContentLoaded', function() {
    const pubBtn = document.getElementById('pub-button');
    pubBtn.addEventListener('click', publish);
});
