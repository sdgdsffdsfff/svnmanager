define([
'kernel',
'react',
'components/template/Mod'
], function( core, React, Mod ){

    return function( opt ){
        var div = $('<div />', {
            id:'id'+ core.random(10)
        });

        var body = React.render(React.createElement(Mod, null), div[0]);

        console.log( opt.bd.constructor );

        return {
            elem: div,
            body: body
        }
    }
});