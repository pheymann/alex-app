import Header from "./Header";

export default function BasicPage(props) {
  return (
    <div>
      <Header awsFetch={ props.awsFetch }
              language={ props.language }
              setLanguage={ props.setLanguage }
              signOut={ props.signOut } />

      {props.children}
    </div>
  );
}
