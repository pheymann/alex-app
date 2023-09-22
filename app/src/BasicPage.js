import Header from "./Header";

export default function BasicPage(props) {
  return (
    <div>
      <Header awsFetch={ props.awsFetch } signOut={ props.signOut } />

      {props.children}
    </div>
  );
}
