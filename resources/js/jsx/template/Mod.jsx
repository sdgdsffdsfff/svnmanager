define([
'kernel',
'react'
], function( core, React ){
    var Mod = React.createClass({
        render: function(){
            return (
                <div className="mod base-style">
                    <div className="hd">
                        <h4></h4>
                    </div>
                    <div className="bd"></div>
                </div>
            )
        }
    });

    return function(){
        return <Mod />;
    }
});