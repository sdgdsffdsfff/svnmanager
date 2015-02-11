define([
'kernel',
'react'
], function( core, React ){
    var Mod = React.createClass({
        render: function(){
            var titleEl = null;
            if( this.props.title ){
                titleEl = <div className="hd">{this.props.title}</div>
            }

            return (
                <div className="mod">
                    {titleEl}
                    <div className="bd">
                    </div>
                </div>
            )
        }
    });

    return Mod;
});