document.addEventListener("DOMContentLoaded", function () {
    const step1 = document.getElementById("kp-wr-step-1");
    const step2 = document.getElementById("kp-wr-step-2");
    const step3 = document.getElementById("kp-wr-step-3");

    const skipStep2 = document.getElementById("kp-wr-skip2");

    const stars = document.querySelectorAll(".kp-star-rating");
    const nextButtonStep1 = document.querySelector("#kp-wr-step-1 button");
    // const nextButtonStep2 = document.getElementById("kp-wr-btn-2");
    const fileInput = document.getElementById("dropzone-file");
    const doneButtonStep3 = document.querySelector("#kp-wr-step-3 button");
 
    var preview = document.getElementById('preview-image');

    let closeModelKPWR = document.getElementById('Writereview-modal-btn')

    let starRating = 0;
    let fileData = null;
    let textareaData = '';
    let userData = { firstName: '', lastName: '', email: '',message:"",starRating:0 ,fileData:null };

    // Function to show a specific step and hide others
    function showStep(stepToShow) {
      [step1, step2, step3].forEach(step => {
        step.style.display = step === stepToShow ? "block" : "none";
      });
    }

    skipStep2.addEventListener("click",function(){
      showStep(step3);
    })

    // Function to reset all values
    function resetValues() {
      starRating = 0;
      textareaData = '';
      fileData = null;
      userData = { firstName: '', lastName: '', email: '', message:"" , starRating:0 , fileData:null};

      // Reset star ratings
      stars.forEach(star => star.classList.remove('text-yellow-500'));

      // Reset file input (clear selection)
      fileInput.value = null;
      preview.src = "/assets/images/reviewpage/upload-img.svg";

      // Reset form values if needed (e.g., clear textarea)
      document.getElementById("message").value = '';
      document.getElementById("first_name").value = '';
      document.getElementById("last_name").value = '';
      document.getElementById("email").value = '';

      // Show Step 1
      showStep(step1);
    }

    closeModelKPWR.addEventListener("click",function(){
      resetValues()
    })
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
      });
    });
    nextButtonStep1.addEventListener("click", function () {
      textareaData = document.getElementById("message").value;
      if (textareaData.trim().length<=10 || starRating<=0) {
          if(textareaData.trim().length<=10  && starRating<=0){
            alert('Please fill in the star rating and review comments before moving to the next step.');
            return;
          }else{
            if(starRating<=0){
              alert('Please fill in the star rating before moving to the next step.');
            return;
            }
            if(textareaData.trim().length<=10){
              alert('Please fill in review moren then 10 words before moving to the next step.');
            return;
            }
          }
      }else{
        return showStep(step2);
      }
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
    // nextButtonStep2.addEventListener("click", function () {
    //     if(fileData&&Object.keys(fileData).length>0){
    //       return showStep(step3);
    //     }else{
    //       alert('Please select the photo before moving to the next step or click on "I’ll Do It Later".');
    //     }
    // });

    function validateEmail(emailInput) {
      const emailRegex = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/;

      if (!emailRegex.test(emailInput)) {
        alert("Invalid email address. Please enter a valid email.");
        return false
      }else{
        return true
      }
    }

    // Step 3: When done button is clicked, submit data and reset values
    doneButtonStep3.addEventListener("click", function () {
      // Capture user data in Step 4
      userData.firstName = document.getElementById("first_name").value;
      userData.lastName = document.getElementById("last_name").value;
      userData.email = document.getElementById("email").value;
      userData.starRating = starRating
      userData.message = textareaData
      userData.fileData = fileData

      if(userData&&Object.keys(userData).length>0){
        if(userData.firstName.toString().trim().length>0&&userData.lastName.toString().trim().length>0&&userData.email.toString().trim().length>0&&validateEmail(userData.email)){
          console.log("Step 4 - User Data:", userData);
          resetValues();
          closeModelKPWR.click()
          alert("Data submitted!");
        }else{
          if(userData.firstName&&userData.firstName.toString().trim().length===0){
            return alert("Please fill in first name.")
          }else{
            if(userData.lastName&&userData.lastName.toString().trim().length===0){
              return alert("Please fill in last name.")
            }else{
              if(userData.email.toString().trim().length===0){
                return alert("Please fill in first name.")
              }else{
                return validateEmail(userData.email)
              }
            }
          }
        }
      }
    });
});