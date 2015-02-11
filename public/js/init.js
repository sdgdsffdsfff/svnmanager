define(['kernel', 'core/delegate'], function( core, delegate ){
    if(_inlineCodes && _inlineCodes.length){
        $.map(_inlineCodes, function(fn){
            typeof fn === 'function' && fn()
        })
    }

    delegate();
});