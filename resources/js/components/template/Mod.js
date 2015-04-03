define([
'kernel',
'react'
], function( core, React ){
    var Mod = React.createClass({displayName: "Mod",
        render: function(){
            return (
                React.createElement("div", {className: "mod base-style"}, 
                    React.createElement("div", {className: "hd"}, 
                        React.createElement("h4", null)
                    ), 
                    React.createElement("div", {className: "bd"})
                )
            )
        }
    });

    return function(){
        return React.createElement(Mod, null);
    }
});