define([
'kernel',
'react'
], function( core, React ){
    var Mod = React.createClass({displayName: "Mod",
        render: function(){
            var titleEl = null;
            if( this.props.title ){
                titleEl = React.createElement("div", {className: "hd"}, this.props.title)
            }

            return (
                React.createElement("div", {className: "mod"}, 
                    titleEl, 
                    React.createElement("div", {className: "bd"}
                    )
                )
            )
        }
    });

    return Mod;
});