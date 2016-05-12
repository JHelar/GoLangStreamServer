function SelectText(element) {
    var doc = document
        , text = doc.getElementById(element)
        , range, selection
    ;    
    if (doc.body.createTextRange) {
        range = document.body.createTextRange();
        range.moveToElementText(text);
        range.select();
    } else if (window.getSelection) {
        selection = window.getSelection();        
        range = document.createRange();
        range.selectNodeContents(text);
        selection.removeAllRanges();
        selection.addRange(range);
    }
}


$(document).ready(function(){
	$("#note-to-hex-form").submit(function(event){
		event.preventDefault();
		var $form = $(this), 
			term = $form.find("input[name='note']").val(),
			url = $form.attr("action");
		var posting = $.post(url, {note: term}, "text");
		posting.done(function(data){
			var content = data;
			$("#result").empty().append(content);
		});
	});
	$("#copy-button").on("click", function(){
		SelectText("result");
		try{
			var successful = document.execCommand('copy');
			var msg = successful ? 'OK' : 'FAIL';
			console.log('Coppy was ' + msg);	
		}catch(err){
			console.log("Unable to copy" + err);
		}
		
	});
});