define([
'kernel',
'core/delegate',
'react'
],
function( core, delegate, React ){
    if(_inlineCodes && _inlineCodes.length){
        $.map(_inlineCodes, function(fn){
            typeof fn === 'function' && fn()
        })
    }

    React.initializeTouchEvents(true);

    delegate();
});