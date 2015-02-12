define([
'kernel',
'react',
'components/template/Confirm'
],
function(core, React, Confirm){
    return function( opt ){

        var elem = $('<div class="ui-flyout box confirm" />');
        React.render(React.createElement(Confirm, React.__spread({},  opt)), elem[0]);

        return elem;
    }
});