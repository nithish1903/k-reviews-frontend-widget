const stepperCancleBtn = document.getElementById("AskAQuestion-modal-btn");

function showStep() {
    // Check if textarea is filled in step-1
    const textareaValue = document.getElementById('message_aq').value;
    if (textareaValue.trim() === '') {
        alert('Please fill in the textarea before moving to the next step.');
        return;
    }

    // Show step-2 and hide step-1
    document.getElementById('kp-aq-step-1').classList.add('hidden');
    document.getElementById('kp-aq-step-2').classList.remove('hidden');
}

function submitAndLogData() {
    // Get values from step-2
    const firstName = document.getElementById('first_name_aq').value;
    const lastName = document.getElementById('last_name_aq').value;
    const email = document.getElementById('email_aq').value;

    // Validate if all fields in step-2 are filled
    if (firstName.trim() === '' || lastName.trim() === '' || email.trim() === '') {
        alert('Please fill in all fields before submitting.');
        return;
    }

    // You can now submit the data or log it as needed.
    // For now, let's log the data to the console.
    console.log('Textarea:', document.getElementById('message_aq').value);
    console.log('First Name:', firstName);
    console.log('Last Name:', lastName);
    console.log('Email:', email);

    // Reset values and go back to step-1
    document.getElementById('message_aq').value = '';
    document.getElementById('first_name_aq').value = '';
    document.getElementById('last_name_aq').value = '';
    document.getElementById('email_aq').value = '';

    // Show step-1 and hide step-2
    document.getElementById('kp-aq-step-1').classList.remove('hidden');
    document.getElementById('kp-aq-step-2').classList.add('hidden');
    stepperCancleBtn.click()
}