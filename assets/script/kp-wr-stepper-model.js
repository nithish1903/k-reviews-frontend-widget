document.addEventListener("DOMContentLoaded", function () {
    const step1 = document.getElementById("kp-wr-step-1");
    const step2 = document.getElementById("kp-wr-step-2");
    const step3 = document.getElementById("kp-wr-step-3");
    const step4 = document.getElementById("kp-wr-step-4");

    const stepperCancleBtn = document.getElementById("kp-wr-stepper-cancle");

    const stars = document.querySelectorAll(".kp-star-rating");
    const fileInput = document.getElementById("dropzone-file");
    const nextButtonStep3 = document.querySelector("#kp-wr-step-3 button");
    const doneButtonStep4 = document.querySelector("#kp-wr-step-4 button");
    var preview = document.getElementById('preview-image');

    let starRating = 0;
    let fileData = null;
    let textareaData = '';
    let userData = { firstName: '', lastName: '', email: '' };

    // Function to show a specific step and hide others
    function showStep(stepToShow) {
      [step1, step2, step3, step4].forEach(step => {
        step.style.display = step === stepToShow ? "block" : "none";
      });
    }

    stepperCancleBtn.addEventListener("click",function(){
        resetValues()
        document.getElementById('Writereview-modal-btn').click();
    })

    // Function to reset all values
    function resetValues() {
      starRating = 0;
      fileData = null;
      textareaData = '';
      userData = { firstName: '', lastName: '', email: '' };

      // Reset star ratings
      stars.forEach(star => star.classList.remove('text-yellow-500'));

      // Reset file input (clear selection)
      fileInput.value = '';

      // Reset form values if needed (e.g., clear textarea)
      document.getElementById("message").value = '';
      document.getElementById("first_name").value = '';
      document.getElementById("last_name").value = '';
      document.getElementById("email").value = '';

      // Show Step 1
      showStep(step1);
    }

    // Event listeners

    // Step 1: When stars are clicked, move to Step 2
    stars.forEach((star, index) => {
      star.addEventListener("click", function () {
        starRating = index + 1;
        stars.forEach((s, i) => {
          if (i < starRating) {
            s.classList.add('text-yellow-500');
          } else {
            s.classList.remove('text-yellow-500');
          }
        });
        showStep(step2);
      });
    });

      

    // Step 2: When file input changes, move to Step 3
    fileInput.addEventListener("change", function () {
      const file = fileInput.files[0];
      var reader = new FileReader();

      reader.onloadend = function () {
          preview.src = reader.result;
      };

      if (file) {
        fileData = {
          name: file.name,
          size: file.size,
          type: file.type,
        };
        reader.readAsDataURL(file);
      }else {
          preview.src = "/assets/images/reviewpage/upload-img.svg";
      }
      showStep(step3);
    });

    // Step 3: When next button is clicked, move to Step 4
    nextButtonStep3.addEventListener("click", function () {
      textareaData = document.getElementById("message").value;
      showStep(step4);
    });

    // Step 4: When done button is clicked, submit data and reset values
    doneButtonStep4.addEventListener("click", function () {
      // Capture user data in Step 4
      userData.firstName = document.getElementById("first_name").value;
      userData.lastName = document.getElementById("last_name").value;
      userData.email = document.getElementById("email").value;

      // Log all data in the console
      console.log("Step 1 - Star Rating:", starRating);
      console.log("Step 2 - File Data:", fileData);
      console.log("Step 3 - Textarea Data:", textareaData);
      console.log("Step 4 - User Data:", userData);

      // Reset values and show Step 1
      resetValues();
      document.getElementById('Writereview-modal-btn').click();
      // Implement your logic to submit data
      alert("Data submitted!");
    });
});