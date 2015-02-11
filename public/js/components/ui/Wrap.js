define(['kernel', 'react'], function( core, React ){
    return function( ui, options ){
        var div = $('<div />', {
            id: 'id'+core.random(10)
        });

        React.render( React.createElement(ui, options), div[0] );
        return div;
    }
});