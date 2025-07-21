package ch.rasc.javersdemo.service;

import java.util.Collection;
import java.util.Collections;

import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;

import ch.rasc.javersdemo.entity.User;

@Service
public class CustomUserDetailsService implements UserDetailsService {

  private final UserService userService;

  public CustomUserDetailsService(UserService userService) {
    this.userService = userService;
  }

  @Override
  public UserDetails loadUserByUsername(String username)
      throws UsernameNotFoundException {
    User user = this.userService.findByUsername(username);
    if (user == null) {
      throw new UsernameNotFoundException("User not found: " + username);
    }

    return new org.springframework.security.core.userdetails.User(user.getUsername(),
        user.getPassword(), user.getEnabled(), true, // accountNonExpired
        true, // credentialsNonExpired
        true, // accountNonLocked
        getAuthorities(user));
  }

  private static Collection<? extends GrantedAuthority> getAuthorities(User user) {
    return Collections
        .singletonList(new SimpleGrantedAuthority("ROLE_" + user.getRole().name()));
  }
}
